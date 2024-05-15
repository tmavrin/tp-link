/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmavrin/tp-link/client"
	"github.com/tmavrin/tp-link/pkg/sms"
)

var smsCmd = &cobra.Command{
	Use:   "sms",
	Short: "SMS commands if SIM card is inserted",
	Long:  `Will not work if no SIM card`,
}

func init() {
	rootCmd.AddCommand(smsCmd)
	smsCmd.AddCommand(smsListCmd)
	smsCmd.AddCommand(smsSendCmd)
	smsCmd.AddCommand(smsFindCmd)

	smsSendCmd.PersistentFlags().String("to", "", "phone number to send the message to")
	smsSendCmd.PersistentFlags().String("content", "", "content of the message")

	smsFindCmd.PersistentFlags().String("from", "", "phone number message was received from")
	smsFindCmd.PersistentFlags().String("phrase", "", "phrase to look for in a message")
}

var smsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SMS Inbox",
	Long:  `Pagination not implemented`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.Authenticate(
			viper.GetString("host"),
			viper.GetString("username"),
			viper.GetString("password"),
		)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		defer c.Close()

		inbox, err := sms.New(c).GetInbox()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		for _, msg := range inbox {
			fmt.Println(msg.String())
		}
	},
}

var smsSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send SMS ",
	Long:  `Send SMS`,
	Run: func(cmd *cobra.Command, args []string) {
		to, err := cmd.Flags().GetString("to")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		content, err := cmd.Flags().GetString("content")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if to == "" || content == "" {
			fmt.Println("--to and --content must be set")
			os.Exit(1)
		}

		c, err := client.Authenticate(
			viper.GetString("host"),
			viper.GetString("username"),
			viper.GetString("password"),
		)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		defer c.Close()

		err = sms.New(c).SendSMS(sms.SMS{
			To:      to,
			Content: content,
		})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		fmt.Println("sms sent successfully")
	},
}

var smsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Searches inbox for keyword from sender",
	Long:  `Search inbox from sender number with message content keywords`,
	Run: func(cmd *cobra.Command, args []string) {
		from, err := cmd.Flags().GetString("from")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		phrase, err := cmd.Flags().GetString("phrase")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if phrase == "" {
			fmt.Println("--phrase must be set")
			os.Exit(1)
		}

		if from == "" {
			fmt.Println("--from not set, searching all messages")
		}

		c, err := client.Authenticate(
			viper.GetString("host"),
			viper.GetString("username"),
			viper.GetString("password"),
		)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		defer c.Close()

		inbox, err := sms.New(c).GetInbox()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		for _, msg := range inbox {
			f := from == "" || from == msg.From
			if f && strings.Contains(msg.Content, phrase) {
				fmt.Printf("Found phrase \"%s\" in message:\n%s\nfrom: %s\n", phrase, msg.Content, msg.From)
				os.Exit(0)
			}
		}

		fmt.Printf("Phrase \"%s\" not found in any message\n", phrase)
		os.Exit(1)
	},
}
