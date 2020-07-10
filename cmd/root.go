package cmd

import (
	"context"
	"fmt"
	"github.com/autom8ter/getter-render/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	dataFile string
	output   string
	source   string
)

var rootCmd = &cobra.Command{
	Use: "getter-render",
	Long: `getter-render extends Hashicorps go-getter library/cli by adding template rendering functionality.
A data file is used to render files fetched from remote sources using go-getter and the go templating language.
`,
	Run: func(cmd *cobra.Command, args []string) {
		renderer := render.NewRenderer()
		if source == "" {
			source = viper.GetString("source")
		}
		if output == "" {
			source = viper.GetString("output")
		}
		if source == "" || output == "" {
			log.Fatal("source(-s --source) and output(-o --output) are required flags")
			return
		}
		if err := renderer.LoadSources(context.Background(), map[string]string{
			output: source,
		}); err != nil {
			log.Fatalf("failed to load source: %s output: %s error: %s", source, output, err.Error())
		}
		if err := renderer.Compile(viper.AllSettings()); err != nil {
			log.Fatalf("failed to compile source: %s output: %s error: %s", source, output, err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output directory(required)")
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "remote file source(required)")
	rootCmd.PersistentFlags().StringVarP(&dataFile, "data", "d", "", "path to data file to render files with (ex data.json)")
}

// initConfig reads in the data file if it exists and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	if dataFile != "" {
		viper.SetConfigFile(dataFile)

		if err := viper.ReadInConfig(); err == nil {
			log.Printf("loaded %s\n", dataFile)
		}
	}
}
