package reconcile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
)

type Report struct {
	NASItems            int   `json:"nas_items"`
	NASBytes            int64 `json:"nas_bytes"`
	GoogleAcceptedItems int   `json:"google_accepted_items"`
	GoogleAcceptedBytes int64 `json:"google_accepted_source_bytes"`
	VerifiedItems       int   `json:"verified_items"`
	VerifiedBytes       int64 `json:"verified_bytes"`
	PendingItems        int   `json:"pending_items"`
	PendingBytes        int64 `json:"pending_bytes"`
	MismatchedItems     int   `json:"mismatched_items"`
	OrphanedReceipts    int   `json:"orphaned_receipts"`
}

func Build(cfg *config.Config, tracker *filetracker.FileTracker) (Report, error) {
	receipts, err := tracker.Receipts()
	if err != nil {
		return Report{}, err
	}
	report := Report{}
	for _, receipt := range receipts {
		if receipt.Version >= 2 && receipt.MediaItemID != "" {
			report.GoogleAcceptedItems++
			report.GoogleAcceptedBytes += receipt.Size
		}
	}
	seen := make(map[string]bool)
	for _, job := range cfg.Jobs {
		matcher, err := filter.Compile(job.IncludePatterns, job.ExcludePatterns)
		if err != nil {
			return report, err
		}
		err = filepath.Walk(job.SourceFolder, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return nil
			}
			if info.IsDir() {
				if path != job.SourceFolder && matcher.IsExcluded(filepath.ToSlash(mustRelative(job.SourceFolder, path))) {
					return filepath.SkipDir
				}
				return nil
			}
			rel := filepath.ToSlash(mustRelative(job.SourceFolder, path))
			if !matcher.IsAllowed(rel) || seen[path] {
				return nil
			}
			seen[path] = true
			report.NASItems++
			report.NASBytes += info.Size()
			receipt, ok := receipts[path]
			if !ok || receipt.Version < 2 || receipt.MediaItemID == "" {
				report.PendingItems++
				report.PendingBytes += info.Size()
				return nil
			}
			hash, err := (filetracker.SHA256Hasher{}).Hash(path)
			if err != nil {
				return err
			}
			if receipt.Size != info.Size() || receipt.Hash != hash {
				report.MismatchedItems++
				return nil
			}
			report.VerifiedItems++
			report.VerifiedBytes += info.Size()
			return nil
		})
		if err != nil {
			return report, fmt.Errorf("scanning %s: %w", job.SourceFolder, err)
		}
	}
	for path, receipt := range receipts {
		if receipt.Version >= 2 && receipt.MediaItemID != "" && !seen[path] {
			report.OrphanedReceipts++
		}
	}
	return report, nil
}

func mustRelative(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}
