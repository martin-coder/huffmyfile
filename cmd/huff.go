/*
Copyright © 2023 Martin Coder <martincoder1@gmail.com>

Use of this source code is governed by an MIT-style
license that can be found in the LICENSE file or at
https://opensource.org/licenses/MIT.
*/

package cmd

import (
	huffmyfile "github.com/martin-coder/huffmyfile/pkg"

	"github.com/spf13/cobra"
)

// huffCmd represents the huff command
var huffCmd = &cobra.Command{
	Use:   "huff",
	Short: "Compresses .txt files into .huff files.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("huff called")
		e := huffmyfile.Encoder{}
		e.EncodeToDefaultOutputFile(args[0])
	},
}

func init() {
	rootCmd.AddCommand(huffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// huffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// huffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
