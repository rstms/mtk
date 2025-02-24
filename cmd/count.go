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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var lsCmd = &cobra.Command{
	Use:   "count",
	Short: "count inbox messages",
	Long: `count inbox messages"
`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := imapLogin()
		cobra.CheckErr(err)

		defer imapLogout(c)

		// List mailboxes
		/*
			mailboxes := make(chan *imap.MailboxInfo, 10)

			done := make(chan error, 1)
			go func() {
				done <- c.List("", "*", mailboxes)
			}()

			log.Println("Mailboxes:")
			for m := range mailboxes {
				log.Println("* " + m.Name)
			}

			if err := <-done; err != nil {
				cobra.CheckErr(err)
			}
		*/

		// Select INBOX
		mbox, err := c.Select("INBOX", false)
		if err != nil {
			cobra.CheckErr(err)
		}
		if viper.GetBool("verbose") {
			log.Println("Flags for INBOX:", mbox.Flags)
		}

		fmt.Printf("%d\n", mbox.Messages)
	},
}

func init() {
	imapCmd.AddCommand(lsCmd)
}
