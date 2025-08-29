package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zallarak/db/cli/internal/colors"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "dbx",
	Short: colors.Gray("db.xyz CLI - Postgres-as-a-Service management tool"),
	Long: colors.Gray("dbx") + colors.White(" is the command line interface for ") + colors.Cyan("db.xyz") + colors.White(", a Postgres-as-a-Service platform.\n\n") +
		colors.White("Manage your organizations, projects, and database instances from the command line."),
}

func Execute() {
	// Disable usage on error for cleaner error messages
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dbx.yaml)")
	rootCmd.PersistentFlags().String("api-url", "http://127.0.0.1:8081", "API base URL")
	rootCmd.PersistentFlags().String("token", "", "API token for authentication")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format: table, json")

	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dbx")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// Silent config loading for minimal output
	}
}