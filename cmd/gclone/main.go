package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gclone/internal/app"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	opts := app.Options{}
	application := app.New()

	cmd := &cobra.Command{
		Use:   "gclone <repo_url>",
		Short: "Clone git repositories using SSH profile aliases",
		Long:  "gclone helps you clone git repositories using a selected SSH host alias from your gclone config.",
		Args:  cobra.ExactArgs(1),
		Example: "  gclone git@gitlab.com:stream-flow/backend/streamflow-api.git --profile profile1\n" +
			"  gclone https://gitlab.com/stream-flow/backend/streamflow-api.git --profile profile1",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RepoURL = args[0]
			opts.Dest, _ = cmd.Flags().GetString("dest")
			opts.Profile, _ = cmd.Flags().GetString("profile")
			opts.DryRun, _ = cmd.Flags().GetBool("dry-run")
			opts.Verbose, _ = cmd.Flags().GetBool("verbose")

			code, err := application.Run(opts)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				os.Exit(code)
			}
			return nil
		},
	}

	cmd.Flags().StringP("profile", "p", "", "Profile name from gclone config")
	cmd.Flags().String("dest", "", "Destination path for the clone")
	cmd.Flags().Bool("dry-run", false, "Print the computed URL and git command without cloning")
	cmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")

	cmd.Version = version
	cmd.SetVersionTemplate(fmt.Sprintf("gclone %s\ncommit: %s\ndate: %s\n", version, commit, date))

	return cmd
}
