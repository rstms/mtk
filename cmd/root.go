/*
Copyright © 2025 Matt Krueger <mkrueger@rstms.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"log"

	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const Version = "0.0.4"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mtk",
	Short: "mail toolkit",
	Long: `mail utilitiy commands
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mtk.yaml)")

	rootCmd.PersistentFlags().StringP("imap-hostname", "i", "", "imap hostname")
	viper.BindPFlag("imap_hostname", rootCmd.PersistentFlags().Lookup("imap-hostname"))

	rootCmd.PersistentFlags().IntP("imap-port", "P", 993, "imap port")
	viper.BindPFlag("imap_port", rootCmd.PersistentFlags().Lookup("imap-port"))

	rootCmd.PersistentFlags().StringP("username", "u", "", "username")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))

	rootCmd.PersistentFlags().StringP("password", "p", "", "password")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose diagnostics")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mtk" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mtk")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

func imapLogin() (*client.Client, error) {

	if viper.GetBool("verbose") {
		log.Println("Connecting to server...")
	}

	// Connect to server
	connectString := fmt.Sprintf("%s:%d", viper.GetString("imap_hostname"), viper.GetInt("imap_port"))
	c, err := client.DialTLS(connectString, nil)
	if err != nil {
		return nil, err
	}

	if viper.GetBool("verbose") {
		log.Println("Connected")
	}

	// Login
	if err := c.Login(viper.GetString("username"), viper.GetString("password")); err != nil {
		return nil, err
	}

	if viper.GetBool("verbose") {
		log.Println("Logged in")
	}

	return c, nil
}

func imapLogout(c *client.Client) {
	// Don't forget to logout
	defer c.Logout()
}
