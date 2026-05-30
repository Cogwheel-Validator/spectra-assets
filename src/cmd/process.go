package cmd

import (
	"spectra-assets/src/config"
	"spectra-assets/src/transformer"

	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process all chains: fetch on-chain validators (BFT) and data-dir entries, then write assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		configs, err := config.LoadConfigs(dataPath)
		if err != nil {
			return err
		}
		return transformer.ProcessAll(assetsPath, configs)
	},
}

func init() {
	rootCmd.AddCommand(processCmd)
}
