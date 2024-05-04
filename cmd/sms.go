/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
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

	smsSendCmd.PersistentFlags().String("to", "", "phone number to send the message to")
	smsSendCmd.PersistentFlags().String("content", "", "content of the message")
}

var smsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SMS Inbox",
	Long:  `Pagination not implemented`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.Authenticate("http://192.168.1.1", "admin", "admin")
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

		c, err := client.Authenticate("http://192.168.1.1", "admin", "admin")
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
