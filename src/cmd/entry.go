package cmd

import (
	"fmt"
	"spectra-assets/src/config"
	"spectra-assets/src/transformer"

	"github.com/spf13/cobra"
)

var entryCmd = &cobra.Command{
	Use:   "process-entry <chain-path> <address>",
	Short: "Process a single data-dir validator entry without hitting the chain API",
	Long: `Validates and processes exactly one validator from the data directory.
No chain API call is made — only the local data/<chain-path>/<address>/ contents
are used. This is the command to run in pull-request CI workflows.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		chainPath, address := args[0], args[1]

		configs, err := config.LoadConfigs(dataPath)
		if err != nil {
			return err
		}

		for _, cfg := range configs {
			if cfg.ChainPath == chainPath {
				return transformer.ProcessEntry(assetsPath, cfg, address)
			}
		}
		return fmt.Errorf("no chain config found for chain-path %q (is there a chain.json under %s/%s?)", chainPath, dataPath, chainPath)
	},
}

func init() {
	rootCmd.AddCommand(entryCmd)
}
