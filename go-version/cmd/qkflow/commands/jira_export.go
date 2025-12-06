package commands

import (
	"path/filepath"

	"github.com/Wangggym/quick-workflow/internal/jira"
	"github.com/spf13/cobra"
)

var (
	exportWithImages bool
	exportOutputDir  string
)

var jiraExportCmd = &cobra.Command{
	Use:   "export <ISSUE-KEY>",
	Short: "Export a Jira issue to files",
	Long: `Export a Jira issue to local files in Markdown format.

By default, exports to /tmp/qkflow/jira/<ISSUE-KEY>/
Use --with-images to download all attachments and images.

The export creates:
  - content.md: Main content in Markdown format
  - README.md: How to use in Cursor
  - attachments/: Downloaded files (if --with-images)

Examples:
  qkflow jira export NA-9245                    # Export text only
  qkflow jira export NA-9245 --with-images      # Export with all images
  qkflow jira export NA-9245 -o ~/exports/      # Custom output directory`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := args[0]

		// Validate issue key
		if !jira.ValidateIssueKey(issueKey) {
			log.Error("Invalid issue key: %s", issueKey)
			return
		}

		// Create Jira client
		client, err := jira.NewClient()
		if err != nil {
			log.Error("Failed to create Jira client: %v", err)
			return
		}

		// Create exporter
		exporter := jira.NewExporter(client)

		// Prepare options
		opts := jira.ExportOptions{
			IssueKey:   issueKey,
			WithImages: exportWithImages,
			OutputDir:  exportOutputDir,
		}

		// Show progress
		log.Info("🔍 Fetching %s from Jira...", issueKey)

		// Export
		result, err := exporter.Export(opts)
		if err != nil {
			log.Error("Failed to export: %v", err)
			return
		}

		// Display results
		log.Info("")
		log.Success("Jira issue %s exported successfully!", issueKey)
		log.Info("")
		log.Info("📁 Location: %s/", result.ExportPath)
		log.Info("📄 Main content: %s", result.ContentFile)

		if exportWithImages && len(result.ImageFiles) > 0 {
			log.Info("🖼️  Images: %d files downloaded", len(result.ImageFiles))
			log.Info("")
			log.Info("Downloaded files:")
			for _, imgPath := range result.ImageFiles {
				log.Info("  - %s", filepath.Base(imgPath))
			}
		}

		// Display Cursor instructions
		log.Info("")
		log.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Info("💡 How to use in Cursor:")
		log.Info("")
		log.Info("1. In Cursor, attach the exported content:")
		log.Info("   @%s", result.ContentFile)
		log.Info("")

		if exportWithImages && len(result.ImageFiles) > 0 {
			log.Info("2. To include images, tell Cursor:")
			log.Info("   \"Please read the images in %s/attachments/\"", result.ExportPath)
			log.Info("")
			log.Info("3. Or simply tell Cursor:")
			log.Info("   \"Read the exported Jira ticket at %s/\"", result.ExportPath)
			log.Info("")
		} else {
			log.Info("2. Or simply tell Cursor:")
			log.Info("   \"Read %s\"", result.ContentFile)
			log.Info("")
		}

		log.Info("4. When done, clean up:")
		log.Info("   qkflow jira clean %s", issueKey)
		log.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Info("")

		// Show hint if images available but not downloaded
		if !exportWithImages {
			hasImages, _ := exporter.HasImages(issueKey)
			if hasImages {
				log.Info("💡 Tip: This issue has attachments. Use --with-images to download them.")
			}
		}
	},
}

func init() {
	jiraExportCmd.Flags().BoolVarP(&exportWithImages, "with-images", "i", false, "Download all attachments and images")
	jiraExportCmd.Flags().StringVarP(&exportOutputDir, "output-dir", "o", "", "Output directory (default: /tmp/qkflow/jira/)")
}
