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
	sourceMap map[string]string
)

var rootCmd = &cobra.Command{
	Use: "getter-render",
	Long: `getter-render extends Hashicorps go-getter library/cli by adding template rendering functionality.
A render file is used to render files fetched from remote sources using go-getter and the go templating language.
`,
	Run: func(cmd *cobra.Command, args []string) {
		renderer := render.NewRenderer()
		if len(sourceMap) == 0 {
			sourceMap = viper.GetStringMapString("source_map")
		}
		if len(sourceMap) == 0 {
			log.Fatal("please add at least one dest:source to `source_map` in render.yaml")
		}
		if err := renderer.LoadSources(context.Background(), sourceMap); err != nil {
			log.Fatalf("failed to load sourceMap: %v error: %s", sourceMap, err.Error())
		}
		if err := renderer.Compile(viper.AllSettings()); err != nil {
			log.Fatalf("failed to compile sourceMap: %v error: %s", sourceMap, err.Error())
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
	rootCmd.PersistentFlags().StringToStringVarP(&sourceMap, "source", "s", map[string]string{}, "source mapping dest: source")
}

// initConfig reads in render file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile("render.yaml")
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		log.Println("loaded render.yaml")
	}
}
