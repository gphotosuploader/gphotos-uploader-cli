package verify

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/reconcile"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewCommand(global *flags.GlobalFlags) *cobra.Command {
	return &cobra.Command{Use: "verify", Short: "Reconcile NAS bytes with durable Google upload receipts", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		cfg, err := config.FromFile(afero.NewOsFs(), filepath.Join(global.CfgDir, app.DefaultConfigFilename), log.Discard)
		if err != nil {
			return err
		}
		repo, err := filetracker.NewLevelDBRepository(filepath.Join(global.CfgDir, "uploaded_files"))
		if err != nil {
			return err
		}
		tracker := filetracker.New(repo)
		defer tracker.Close()
		report, err := reconcile.Build(cfg, tracker)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			return err
		}
		if report.PendingItems > 0 || report.MismatchedItems > 0 || report.OrphanedReceipts > 0 || report.VerifiedItems != report.NASItems {
			return fmt.Errorf("verification incomplete")
		}
		return nil
	}}
}
