package cmd

import "github.com/spf13/cobra"

var (
	dataPath   string
	assetsPath string
)

var rootCmd = &cobra.Command{
	Use:   "spectra-assets",
	Short: "Generate validator logo and config assets for the Spectra explorer",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataPath, "data", "data", "path to the data directory")
	rootCmd.PersistentFlags().StringVar(&assetsPath, "assets", "assets", "path to the assets output directory")
}
