package parse

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"

	parseblocks "github.com/forbole/juno/v3/cmd/parse/blocks"
	parsegenesis "github.com/forbole/juno/v3/cmd/parse/genesis"

	parseprofiles "github.com/desmos-labs/djuno/v2/cmd/parse/profiles"
)

// NewParseCmd returns the Cobra command allowing to parse some data without having to re-sync the whole database
func NewParseCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "parse",
		Short:             "Apply some fixes without the need to re-syncing the whole database from scratch",
		PersistentPreRunE: runPersistentPreRuns(parsecmdtypes.ReadConfigPreRunE(parseCfg)),
	}

	cmd.AddCommand(
		parsegenesis.NewGenesisCmd(parseCfg),
		parseblocks.NewBlocksCmd(parseCfg),
		parseprofiles.NewProfilesCmd(parseCfg),
	)

	return cmd
}

func runPersistentPreRuns(preRun func(_ *cobra.Command, _ []string) error) func(_ *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if root := cmd.Root(); root != nil {
			if root.PersistentPreRunE != nil {
				err := root.PersistentPreRunE(root, args)
				if err != nil {
					return err
				}
			}
		}

		return preRun(cmd, args)
	}
}