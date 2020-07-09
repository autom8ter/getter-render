package cmd

import (
	"context"
	"fmt"
	"github.com/autom8ter/getter-tmpl/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	valuesFile string
)

var rootCmd = &cobra.Command{
	Use: "getter-render",
	Long: `getter-render extends Hashicorps go-getter library/cli by adding template rendering functionality.
A values file is used to render files fetched from remote sources using go-getter and the go templating language.
`,
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := render.NewFileSet()
		sources := viper.GetStringSlice("sources")
		if err := tmpl.Load(context.Background(), sources); err != nil {
			log.Fatalf("failed to load sources: %v error: %s", sources, err.Error())
		}
		if err := tmpl.Compile(viper.AllSettings()); err != nil {
			log.Fatalf("failed to compile sources: %v error: %s", sources, err.Error())
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
	rootCmd.PersistentFlags().StringVarP(&valuesFile, "values", "v", "values.yaml", "values file to render files")
}

// initConfig reads in values file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(valuesFile)
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read in values file: %s", valuesFile)
	}
}
