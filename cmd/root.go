package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "tp-link",
	Short: "A CLI Client for TP-Link Router",
	Long: `Currently implements only SMS Inbox list and send commands.
	Tested only on Archer MR600`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tp-link.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		fmt.Printf("Using config file: %s\n", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tp-link")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config not found. Creating config with defaults at $HOME/.tp-link.yaml ")
			viper.Set("username", "admin")
			viper.Set("password", "admin")
			viper.Set("host", "http://192.168.1.1")
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			err = viper.WriteConfigAs(fmt.Sprintf("%s/.tp-link.yaml", home))
			cobra.CheckErr(err)
		}
	}
}
