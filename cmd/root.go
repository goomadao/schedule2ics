/*
Copyright © 2020 MADAO <goomadao@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/goomadao/schedule2ics/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, scheduleFile, outFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "schedule2ics",
	Short: "A tool to transfer your schedule into .ics",
	Long: `A tool to transfer your schedule into .ics
You can import the .ics file into your Google Calendar or Apple Calendar.`,

	Run: func(cmd *cobra.Command, args []string) {
		util.Classes2ICS(util.Xlsx2Classes(scheduleFile), outFile)
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./schedule2ics.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outFile, "output", "o", "schedule.ics", "output .ics file")
	rootCmd.PersistentFlags().StringVarP(&scheduleFile, "schedule", "s", "", "schedule file(must be .xls or .xlsx)")
	rootCmd.MarkPersistentFlagFilename("schedule")
	rootCmd.MarkPersistentFlagRequired("schedule")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("schedule2ics")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
