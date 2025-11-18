package jira

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Cleaner handles cleaning up exported Jira files
type Cleaner struct {
	baseDir string
}

// NewCleaner creates a new cleaner
func NewCleaner() *Cleaner {
	return &Cleaner{
		baseDir: filepath.Join("/tmp", "qkflow", "jira"),
	}
}

// CleanOptions contains options for cleaning
type CleanOptions struct {
	IssueKey string
	All      bool
	DryRun   bool
}

// CleanResult contains the result of a clean operation
type CleanResult struct {
	IssueKey   string
	Path       string
	Size       int64
	FileCount  int
	Deleted    bool
	Error      error
}

// Clean cleans up exported Jira files
func (c *Cleaner) Clean(opts CleanOptions) ([]CleanResult, error) {
	if opts.All {
		return c.cleanAll(opts.DryRun)
	}
	
	if opts.IssueKey == "" {
		return nil, fmt.Errorf("issue key is required")
	}

	result, err := c.cleanIssue(opts.IssueKey, opts.DryRun)
	if err != nil {
		return nil, err
	}

	return []CleanResult{*result}, nil
}

// cleanIssue cleans up a single issue's export
func (c *Cleaner) cleanIssue(issueKey string, dryRun bool) (*CleanResult, error) {
	issueDir := filepath.Join(c.baseDir, issueKey)

	// Check if directory exists
	info, err := os.Stat(issueDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("export not found for %s", issueKey)
		}
		return nil, fmt.Errorf("failed to stat directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", issueDir)
	}

	// Calculate size and file count
	size, count, err := calculateDirSize(issueDir)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate directory size: %w", err)
	}

	result := &CleanResult{
		IssueKey:  issueKey,
		Path:      issueDir,
		Size:      size,
		FileCount: count,
	}

	// Delete if not dry run
	if !dryRun {
		if err := os.RemoveAll(issueDir); err != nil {
			result.Error = err
			return result, fmt.Errorf("failed to delete directory: %w", err)
		}
		result.Deleted = true
	}

	return result, nil
}

// cleanAll cleans up all exported issues
func (c *Cleaner) cleanAll(dryRun bool) ([]CleanResult, error) {
	// Check if base directory exists
	if _, err := os.Stat(c.baseDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no exports found")
	}

	// List all issue directories
	entries, err := os.ReadDir(c.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var results []CleanResult
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		issueKey := entry.Name()
		result, err := c.cleanIssue(issueKey, dryRun)
		if err != nil {
			// Log error but continue
			results = append(results, CleanResult{
				IssueKey: issueKey,
				Error:    err,
			})
			continue
		}
		results = append(results, *result)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no exports found")
	}

	return results, nil
}

// ListExports lists all exported issues
func (c *Cleaner) ListExports() ([]ExportInfo, error) {
	// Check if base directory exists
	if _, err := os.Stat(c.baseDir); os.IsNotExist(err) {
		return nil, nil // No exports is not an error
	}

	// List all issue directories
	entries, err := os.ReadDir(c.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var exports []ExportInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		issueKey := entry.Name()
		issueDir := filepath.Join(c.baseDir, issueKey)

		size, count, err := calculateDirSize(issueDir)
		if err != nil {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		exports = append(exports, ExportInfo{
			IssueKey:  issueKey,
			Path:      issueDir,
			Size:      size,
			FileCount: count,
			ModTime:   info.ModTime(),
		})
	}

	return exports, nil
}

// ExportInfo contains information about an export
type ExportInfo struct {
	IssueKey  string
	Path      string
	Size      int64
	FileCount int
	ModTime   interface{}
}

// Helper functions

func calculateDirSize(path string) (int64, int, error) {
	var size int64
	var count int

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
			count++
		}
		return nil
	})

	return size, count, err
}

// FormatCleanResult formats a clean result for display
func FormatCleanResult(result CleanResult, dryRun bool) string {
	var sb strings.Builder

	if result.Error != nil {
		sb.WriteString(fmt.Sprintf("❌ Error cleaning %s: %v\n", result.IssueKey, result.Error))
		return sb.String()
	}

	if dryRun {
		sb.WriteString(fmt.Sprintf("🗑️  Would delete: %s\n", result.IssueKey))
	} else {
		sb.WriteString(fmt.Sprintf("✅ Cleaned: %s\n", result.IssueKey))
	}
	sb.WriteString(fmt.Sprintf("   Path: %s\n", result.Path))
	sb.WriteString(fmt.Sprintf("   Size: %s (%d files)\n", formatSize(result.Size), result.FileCount))

	return sb.String()
}

// FormatExportsList formats a list of exports for display
func FormatExportsList(exports []ExportInfo) string {
	var sb strings.Builder

	sb.WriteString("📦 Exported Jira Issues:\n\n")
	
	if len(exports) == 0 {
		sb.WriteString("No exports found.\n")
		return sb.String()
	}

	totalSize := int64(0)
	for _, exp := range exports {
		sb.WriteString(fmt.Sprintf("  • %s\n", exp.IssueKey))
		sb.WriteString(fmt.Sprintf("    Path: %s\n", exp.Path))
		sb.WriteString(fmt.Sprintf("    Size: %s (%d files)\n", formatSize(exp.Size), exp.FileCount))
		sb.WriteString("\n")
		totalSize += exp.Size
	}

	sb.WriteString(fmt.Sprintf("Total: %d exports, %s\n", len(exports), formatSize(totalSize)))

	return sb.String()
}

