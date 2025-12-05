package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Wangggym/quick-workflow/internal/config"
	"github.com/Wangggym/quick-workflow/internal/github"
	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/Wangggym/quick-workflow/internal/logger"
	"github.com/Wangggym/quick-workflow/internal/utils"
	"github.com/Wangggym/quick-workflow/internal/watcher"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch daemon for auto-updating Jira",
	Long:  `Watch daemon monitors your PRs and automatically updates Jira when merged.`,
}

var watchCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Manually check PRs once",
	Long:  `Manually check for merged PRs and update Jira status.`,
	Run:   runWatchCheck,
}

var watchStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start watch daemon",
	Long:  `Start the watch daemon to monitor PRs in the background.`,
	Run:   runWatchStart,
}

var watchStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop watch daemon",
	Long:  `Stop the running watch daemon.`,
	Run:   runWatchStop,
}

var watchRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart watch daemon",
	Long:  `Restart the watch daemon.`,
	Run:   runWatchRestart,
}

var watchStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show watch daemon status",
	Long:  `Show the current status of the watch daemon.`,
	Run:   runWatchStatus,
}

var watchDaemonCmd = &cobra.Command{
	Use:    "daemon",
	Short:  "Run as daemon (internal use)",
	Long:   `Run as daemon process. This command is for internal use only.`,
	Hidden: true,
	Run:    runWatchDaemon,
}

var watchInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and start watch daemon with auto-start",
	Long:  `Install watch daemon as a system service (launchd on macOS) with auto-start on login.`,
	Run:   runWatchInstall,
}

var watchUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall watch daemon and remove auto-start",
	Long:  `Uninstall the watch daemon service and remove auto-start configuration.`,
	Run:   runWatchUninstall,
}

var watchLogCmd = &cobra.Command{
	Use:   "log",
	Short: "View watch daemon logs",
	Long:  `View the watch daemon log file.`,
	Run:   runWatchLog,
}

var watchHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View PR processing history",
	Long:  `View the history of processed PRs.`,
	Run:   runWatchHistory,
}

var watchConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show watch configuration",
	Long:  `Show the current watch daemon configuration.`,
	Run:   runWatchConfig,
}

var dryRun bool
var followLog bool
var logLines int
var historyDays int

func init() {
	watchCheckCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate run without actually updating Jira")
	watchLogCmd.Flags().BoolVar(&followLog, "follow", false, "Follow log output (like tail -f)")
	watchLogCmd.Flags().IntVar(&logLines, "last", 50, "Show last N lines")
	watchHistoryCmd.Flags().IntVar(&historyDays, "days", 7, "Show history for last N days")

	watchCmd.AddCommand(watchCheckCmd)
	watchCmd.AddCommand(watchStartCmd)
	watchCmd.AddCommand(watchStopCmd)
	watchCmd.AddCommand(watchRestartCmd)
	watchCmd.AddCommand(watchStatusCmd)
	watchCmd.AddCommand(watchInstallCmd)
	watchCmd.AddCommand(watchUninstallCmd)
	watchCmd.AddCommand(watchLogCmd)
	watchCmd.AddCommand(watchHistoryCmd)
	watchCmd.AddCommand(watchConfigCmd)
	watchCmd.AddCommand(watchDaemonCmd)
}

func runWatchCheck(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		log.Error("Configuration not found. Please run 'qkflow init' first")
		return
	}

	// Validate configuration
	if cfg.GitHubToken == "" {
		log.Error("GitHub token not configured")
		return
	}

	if cfg.JiraServiceAddress == "" || cfg.Email == "" || cfg.JiraAPIToken == "" {
		log.Error("Jira not configured. Please run 'qkflow init' first")
		return
	}

	// Get GitHub client
	ghClient, err := github.NewClient()
	if err != nil {
		log.Error("Failed to create GitHub client: %v", err)
		return
	}

	if dryRun {
		log.Info("🔍 Dry-run mode: No Jira updates will be made")
		log.Info("")
	}

	// Initialize components
	// Level will be automatically loaded from environment variable or default value
	watcherLogger, err := logger.NewLogger(&logger.LoggerOptions{
		Type:     logger.LoggerTypeFile,
		FileName: "watch.log",
		// Level omitted - will use QKFLOW_LOG_LEVEL env var or default LevelInfo
	})
	if err != nil {
		log.Error("Failed to create logger: %v", err)
		return
	}
	defer watcherLogger.Close()

	state, err := watcher.NewState()
	if err != nil {
		log.Error("Failed to load state: %v", err)
		return
	}

	watchingList, err := watcher.NewWatchingList()
	if err != nil {
		log.Error("Failed to load watching list: %v", err)
		return
	}

	log.Info("Checking %d watching PRs...", watchingList.Count())
	log.Info("")

	if watchingList.Count() == 0 {
		log.Info("No PRs in watching list")
		log.Info("PRs will be added automatically when you create them with 'qkflow pr create'")
		return
	}

	jiraClient, err := jira.NewClient()
	if err != nil {
		log.Error("Failed to create Jira client: %v", err)
		return
	}

	statusCache, err := jira.NewStatusCache()
	if err != nil {
		log.Error("Failed to load Jira status cache: %v", err)
		return
	}

	// Create checker and processor
	checker := watcher.NewChecker(ghClient, watcherLogger)
	processor := watcher.NewProcessor(jiraClient, statusCache, watcherLogger)

	// Check for merged PRs from watching list
	mergedPRs, err := checker.CheckMergedPRs(watchingList, state)
	if err != nil {
		log.Error("Failed to check PRs: %v", err)
		watcherLogger.Error("Failed to check PRs: %v", err)
		return
	}

	if len(mergedPRs) == 0 {
		log.Success("✅ No newly merged PRs found")
		watcherLogger.Info("No newly merged PRs found")
		return
	}

	log.Info("Found %d newly merged PR(s) with Jira tickets", len(mergedPRs))
	log.Info("")

	// Process each PR
	for _, pr := range mergedPRs {
		log.Info("📋 PR #%d: %s", pr.Number, pr.Title)
		log.Info("   Branch: %s", pr.Branch)
		log.Info("   Jira: %v", pr.JiraTickets)
		log.Info("   Merged: %s by %s", pr.MergedAt, pr.MergedBy)
		log.Info("")

		if dryRun {
			// Dry-run: just log what would happen
			for _, ticket := range pr.JiraTickets {
				projectKey := watcher.GetProjectFromTicket(ticket)
				mapping, err := statusCache.GetProjectStatus(projectKey)
				if err != nil || mapping == nil {
					log.Warning("   ⚠️  %s: No status mapping configured for project %s", ticket, projectKey)
					continue
				}

				log.Info("   Would update %s → %s", ticket, mapping.PRMergedStatus)
			}
			log.Info("")
			continue
		}

		// Actually process the PR
		processedPR := processor.ProcessMergedPR(pr)

		// Display results
		for _, update := range processedPR.JiraUpdates {
			if update.Success {
				log.Success("   ✅ %s: %s → %s", update.Ticket, update.OldStatus, update.NewStatus)
			} else {
				log.Error("   ❌ %s: %s", update.Ticket, update.Error)
			}
		}

		log.Info("")

		// Save to state
		if err := state.AddProcessedPR(processedPR); err != nil {
			log.Warning("Failed to save processed PR to state: %v", err)
		}

		// Remove from watching list
		for _, watchingPR := range watchingList.GetAll() {
			if watchingPR.PRNumber == pr.Number {
				if err := watchingList.Remove(watchingPR.Owner, watchingPR.Repo, pr.Number); err != nil {
					watcherLogger.Warning("Failed to remove PR #%d from watching list: %v", pr.Number, err)
				} else {
					watcherLogger.Info("Removed PR #%d from watching list", pr.Number)
				}
				break
			}
		}
	}

	// Update last check time
	if err := state.UpdateLastCheckTime(); err != nil {
		watcherLogger.Warning("Failed to update last check time: %v", err)
	}

	// Clean old records
	retentionDays := 7 // Default from config
	if err := state.CleanOldRecords(retentionDays); err != nil {
		watcherLogger.Warning("Failed to clean old records: %v", err)
	}

	// Note: Log cleanup is handled by the logger itself, no manual cleanup needed

	if dryRun {
		log.Info("🔍 Dry-run completed. No changes were made.")
	} else {
		log.Success("✅ Processed %d PR(s)", len(mergedPRs))
	}

	// Log file path is in the config directory
	configDir, _ := utils.GetConfigDir()
	if configDir != "" {
		logPath := filepath.Join(configDir, "watch.log")
		log.Info("\n📝 Logs: %s", logPath)
	}
}

func runWatchDaemon(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		fmt.Fprintf(os.Stderr, "Configuration not found. Please run 'qkflow init' first\n")
		os.Exit(1)
	}

	// Create daemon with default schedule config
	daemon, err := watcher.NewDaemon(cfg, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create daemon: %v\n", err)
		os.Exit(1)
	}

	// Start daemon (blocks until signal)
	if err := daemon.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
		os.Exit(1)
	}
}

func runWatchStart(cmd *cobra.Command, args []string) {
	// Check if already running
	running, pid, err := watcher.IsRunning()
	if err != nil {
		log.Error("Failed to check daemon status: %v", err)
		return
	}

	if running {
		log.Warning("Watch daemon is already running (PID: %d)", pid)
		return
	}

	log.Info("🚀 Starting watch daemon...")

	// Fork daemon process
	execPath, err := os.Executable()
	if err != nil {
		log.Error("Failed to get executable path: %v", err)
		return
	}

	// Start daemon in background
	procAttr := &os.ProcAttr{
		Files: []*os.File{nil, nil, nil}, // Detach from terminal
	}

	process, err := os.StartProcess(execPath, []string{execPath, "watch", "daemon"}, procAttr)
	if err != nil {
		log.Error("Failed to start daemon: %v", err)
		return
	}

	// Release the process
	process.Release()

	// Give it a moment to start
	time.Sleep(time.Second)

	// Verify it's running
	running, pid, _ = watcher.IsRunning()
	if running {
		log.Success("✅ Watch daemon started successfully (PID: %d)", pid)
	} else {
		log.Warning("Daemon may have failed to start. Check logs for details.")
	}
}

func runWatchStop(cmd *cobra.Command, args []string) {
	running, pid, err := watcher.IsRunning()
	if err != nil {
		log.Error("Failed to check daemon status: %v", err)
		return
	}

	if !running {
		log.Info("Watch daemon is not running")
		return
	}

	log.Info("🛑 Stopping watch daemon (PID: %d)...", pid)

	// Find process
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Error("Failed to find process: %v", err)
		return
	}

	// Send SIGTERM
	if err := process.Signal(syscall.SIGTERM); err != nil {
		log.Error("Failed to send stop signal: %v", err)
		return
	}

	// Wait for process to stop (with timeout)
	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		running, _, _ := watcher.IsRunning()
		if !running {
			log.Success("✅ Watch daemon stopped")
			return
		}
	}

	log.Warning("Daemon did not stop within timeout. It may still be shutting down.")
}

func runWatchRestart(cmd *cobra.Command, args []string) {
	log.Info("🔄 Restarting watch daemon...")

	// Stop if running
	running, _, _ := watcher.IsRunning()
	if running {
		runWatchStop(cmd, args)
		time.Sleep(time.Second)
	}

	// Start
	runWatchStart(cmd, args)
}

func runWatchStatus(cmd *cobra.Command, args []string) {
	state, err := watcher.NewState()
	if err != nil {
		log.Error("Failed to load state: %v", err)
		return
	}

	cfg := config.Get()

	log.Info("")
	log.Info("🔍 Watch Daemon Status")
	log.Info("")

	// Check if running
	running, pid, _ := watcher.IsRunning()

	if running {
		uptime := time.Since(state.DaemonStartTime)
		log.Info("Status: Running ✅")
		log.Info("PID: %d", pid)
		log.Info("Started: %s (uptime: %s)",
			state.DaemonStartTime.Format("2006-01-02 15:04:05"),
			formatDuration(uptime))

		if !state.LastCheckTime.IsZero() {
			lastCheck := time.Since(state.LastCheckTime)
			log.Info("Last Check: %s (%s ago)",
				state.LastCheckTime.Format("2006-01-02 15:04:05"),
				formatDuration(lastCheck))

			// Calculate next check
			scheduler := watcher.NewScheduler(nil)
			nextCheck := scheduler.CalculateNextCheckTime(time.Now())
			log.Info("Next Check: %s (%s)",
				nextCheck.Format("2006-01-02 15:04:05"),
				scheduler.FormatNextCheckTime(nextCheck))
		}
	} else {
		log.Info("Status: Stopped ⏸️")
		log.Info("")
		log.Warning("Watch daemon is not running")
	}

	log.Info("")
	log.Info("📊 Statistics (last 7 days):")
	recentPRs := state.GetRecentPRs(7)
	log.Info("  PRs Processed: %d", len(recentPRs))

	successCount := 0
	for _, pr := range recentPRs {
		for _, update := range pr.JiraUpdates {
			if update.Success {
				successCount++
			}
		}
	}
	log.Info("  Jira Updated: %d", successCount)
	log.Info("  Errors: %d", state.Stats.TotalErrors)

	if cfg != nil {
		log.Info("")
		log.Info("📋 Configuration:")
		log.Info("  Repository: %s/%s", cfg.GitHubOwner, cfg.GitHubRepo)
		log.Info("  Author: %s (only your PRs)", cfg.GitHubOwner)
		log.Info("  Check Interval: 15 minutes (daytime)")
		log.Info("  Night Checks: 02:00, 06:00")
		log.Info("  Log Retention: 7 days")
	}

	log.Info("")
	log.Info("📝 Files:")
	if configDir, err := utils.GetConfigDir(); err == nil && configDir != "" {
		logPath := filepath.Join(configDir, "watch.log")
		log.Info("  Log: %s", logPath)
	}

	if cfg != nil {
		configDir, _ := utils.GetConfigDir()
		if configDir != "" {
			log.Info("  State: %s/watch-state.json", configDir)
			log.Info("  Config: %s/config.yaml", configDir)
		}
	}

	log.Info("")
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute"
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minute(s)", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%d hour(s)", hours)
		}
		return fmt.Sprintf("%d hour(s) %d minute(s)", hours, minutes)
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if hours == 0 {
		return fmt.Sprintf("%d day(s)", days)
	}
	return fmt.Sprintf("%d day(s) %d hour(s)", days, hours)
}

func runWatchInstall(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		log.Error("Configuration not found. Please run 'qkflow init' first")
		return
	}

	log.Info("📦 Installing watch daemon...")
	log.Info("")

	// Check prerequisites
	log.Info("Checking prerequisites...")
	if cfg.GitHubToken == "" {
		log.Error("✗ GitHub token not configured")
		return
	}
	log.Info("  ✓ GitHub token configured")

	if cfg.JiraServiceAddress == "" || cfg.Email == "" || cfg.JiraAPIToken == "" {
		log.Error("✗ Jira not configured")
		return
	}
	log.Info("  ✓ Jira configured")

	// Check Jira status mappings
	statusCache, err := jira.NewStatusCache()
	if err != nil {
		log.Error("✗ Failed to load Jira status cache: %v", err)
		return
	}

	mappings, err := statusCache.ListAllMappings()
	if err != nil || len(mappings) == 0 {
		log.Warning("✗ No Jira status mappings found")
		log.Info("  Please run 'qkflow jira setup' to configure status mappings")
		return
	}
	log.Info("  ✓ Jira status mappings found (%d project(s))", len(mappings))

	log.Info("")

	// Check if already installed
	installed, err := watcher.IsLaunchAgentInstalled()
	if err != nil {
		log.Error("Failed to check installation status: %v", err)
		return
	}

	if installed {
		log.Warning("Watch daemon is already installed")
		log.Info("To reinstall, run 'qkflow watch uninstall' first")
		return
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		log.Error("Failed to get executable path: %v", err)
		return
	}

	// Install launch agent
	log.Info("Creating launch agent...")
	if err := watcher.InstallLaunchAgent(execPath); err != nil {
		log.Error("✗ Failed to install launch agent: %v", err)
		return
	}
	log.Info("  ✓ Generated plist file")
	log.Info("  ✓ Installed to ~/Library/LaunchAgents/")

	// Give it a moment to start
	time.Sleep(2 * time.Second)

	// Check if running
	running, pid, _ := watcher.IsRunning()
	if running {
		log.Info("  ✓ Watch daemon started (PID: %d)", pid)
	} else {
		log.Warning("  ⚠️  Daemon may take a moment to start")
	}

	log.Info("")
	log.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Info("")
	log.Success("🎉 Installation complete!")
	log.Info("")
	log.Info("The watch daemon is now running and will:")
	log.Info("  • Check your PRs every 15 minutes (8:30-24:00)")
	log.Info("  • Check twice during night (2:00, 6:00)")
	log.Info("  • Update Jira when PRs merge")
	log.Info("  • Start automatically on login")
	log.Info("")
	log.Info("📊 Useful commands:")
	log.Info("  • Check status: qkflow watch status")
	log.Info("  • View logs: qkflow watch log --follow")
	log.Info("  • View history: qkflow watch history")
	log.Info("  • Stop daemon: qkflow watch stop")
	log.Info("  • Uninstall: qkflow watch uninstall")
	log.Info("")

	configDir, _ := utils.GetConfigDir()
	if configDir != "" {
		logPath := filepath.Join(configDir, "watch.log")
		log.Info("📝 Logs: %s", logPath)
	}
}

func runWatchUninstall(cmd *cobra.Command, args []string) {
	log.Info("🗑️  Uninstalling watch daemon...")
	log.Info("")

	// Check if installed
	installed, err := watcher.IsLaunchAgentInstalled()
	if err != nil {
		log.Error("Failed to check installation status: %v", err)
		return
	}

	if !installed {
		log.Info("Watch daemon is not installed")
		return
	}

	// Stop daemon if running
	running, pid, _ := watcher.IsRunning()
	if running {
		log.Info("Stopping daemon (PID: %d)...", pid)

		process, err := os.FindProcess(pid)
		if err == nil {
			process.Signal(syscall.SIGTERM)
			time.Sleep(time.Second)
		}

		log.Info("  ✓ Daemon stopped")
	}

	// Uninstall launch agent
	log.Info("Removing auto-start...")
	if err := watcher.UninstallLaunchAgent(); err != nil {
		log.Error("✗ Failed to uninstall launch agent: %v", err)
		return
	}

	plistPath, _ := watcher.GetLaunchAgentPath()
	log.Info("  ✓ Removed launch agent")
	if plistPath != "" {
		log.Info("  ✓ Deleted %s", plistPath)
	}

	log.Info("")
	log.Success("✅ Watch daemon completely uninstalled")
	log.Info("")
	log.Info("The daemon will NOT start automatically anymore.")

	configDir, _ := utils.GetConfigDir()
	if configDir != "" {
		log.Info("Logs and history are preserved at ~/.qkflow/")
	}

	log.Info("")
	log.Info("💡 To re-enable later, use 'qkflow watch install'")
}

func runWatchLog(cmd *cobra.Command, args []string) {
	configDir, err := utils.GetConfigDir()
	if err != nil {
		log.Error("Failed to get config directory: %v", err)
		return
	}
	logPath := filepath.Join(configDir, "watch.log")

	if followLog {
		// Follow log (tail -f style)
		log.Info("Following logs: %s", logPath)
		log.Info("Press Ctrl+C to stop")
		log.Info("")

		cmd := exec.Command("tail", "-f", logPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	// Show last N lines
	lines, err := readLastLines(logPath, logLines)
	if err != nil {
		log.Error("Failed to read log file: %v", err)
		return
	}

	log.Info("")
	log.Info("📝 Last %d lines from %s:", logLines, logPath)
	log.Info("")
	for _, line := range lines {
		log.Info(line)
	}
	log.Info("")
}

func runWatchHistory(cmd *cobra.Command, args []string) {
	state, err := watcher.NewState()
	if err != nil {
		log.Error("Failed to load state: %v", err)
		return
	}

	prs := state.GetRecentPRs(historyDays)

	if len(prs) == 0 {
		log.Info("No PRs processed in the last %d days", historyDays)
		return
	}

	log.Info("")
	log.Info("📋 PR Processing History (Last %d days)", historyDays)
	log.Info("")

	for i, pr := range prs {
		if i > 0 {
			log.Info("")
		}

		// Determine overall success
		allSuccess := true
		successCount := 0
		for _, update := range pr.JiraUpdates {
			if update.Success {
				successCount++
			} else {
				allSuccess = false
			}
		}

		status := "✅"
		if !allSuccess {
			if successCount == 0 {
				status = "❌"
			} else {
				status = "⚠️ "
			}
		}

		log.Info("%s PR #%d: %s", status, pr.PRNumber, pr.PRTitle)
		log.Info("   Branch: %s", pr.Branch)
		log.Info("   Merged: %s by %s", pr.MergedAt.Format("2006-01-02 15:04:05"), pr.MergedBy)

		if len(pr.JiraTickets) > 0 {
			tickets := ""
			for j, ticket := range pr.JiraTickets {
				if j > 0 {
					tickets += ", "
				}
				tickets += ticket
			}
			log.Info("   Jira: %s", tickets)
		}

		if len(pr.JiraUpdates) > 0 {
			for _, update := range pr.JiraUpdates {
				if update.Success {
					log.Info("   ✓ %s: %s → %s", update.Ticket, update.OldStatus, update.NewStatus)
				} else {
					log.Info("   ✗ %s: %s", update.Ticket, update.Error)
				}
			}
		}

		log.Info("   Processed: %s", pr.ProcessedAt.Format("2006-01-02 15:04:05"))
	}

	log.Info("")
	log.Info("Total: %d PRs processed, %d Jira updates", len(prs), state.Stats.TotalJiraUpdated)
	log.Info("")
}

func runWatchConfig(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		log.Error("Configuration not found")
		return
	}

	log.Info("")
	log.Info("⚙️  Watch Daemon Configuration")
	log.Info("")

	// Show watching list status
	watchingList, err := watcher.NewWatchingList()
	if err == nil {
		log.Info("📋 Watching List:")
		log.Info("  Total PRs: %d", watchingList.Count())
		if watchingList.Count() > 0 {
			for _, pr := range watchingList.GetAll() {
				log.Info("  • PR #%d: %s/%s", pr.PRNumber, pr.Owner, pr.Repo)
			}
		}
	}

	log.Info("")
	log.Info("⏰ Schedule:")
	log.Info("  Daytime (8:30-24:00): Every 15 minutes")
	log.Info("  Night (0:00-8:30): 2:00, 6:00")
	log.Info("  Log Retention: 7 days")

	log.Info("")
	log.Info("🔔 Notifications:")
	log.Info("  Desktop Notify: Enabled (macOS)")

	// Check installation status
	installed, _ := watcher.IsLaunchAgentInstalled()
	running, pid, _ := watcher.IsRunning()

	log.Info("")
	log.Info("📊 Status:")
	if installed {
		log.Info("  Auto-start: Enabled ✅")
	} else {
		log.Info("  Auto-start: Not installed ❌")
	}

	if running {
		log.Info("  Daemon: Running (PID: %d) ✅", pid)
	} else {
		log.Info("  Daemon: Stopped ⏸️")
	}

	// Show Jira mappings
	statusCache, err := jira.NewStatusCache()
	if err == nil {
		mappings, err := statusCache.ListAllMappings()
		if err == nil && len(mappings) > 0 {
			log.Info("")
			log.Info("🎫 Jira Status Mappings:")
			for _, mapping := range mappings {
				log.Info("  %s:", mapping.ProjectKey)
				log.Info("    PR Created → %s", mapping.PRCreatedStatus)
				log.Info("    PR Merged  → %s", mapping.PRMergedStatus)
			}
		}
	}

	log.Info("")
}

// readLastLines reads the last N lines from a file
func readLastLines(filePath string, n int) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	lines := make([]string, 0)
	currentLine := ""

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = ""
		} else {
			currentLine += string(data[i])
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// Return last n lines
	if len(lines) <= n {
		return lines, nil
	}

	return lines[len(lines)-n:], nil
}
