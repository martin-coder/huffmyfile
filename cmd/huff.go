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
	Short: "Compresses .txt files into .huff files. Usage: `huffmyfile huff [FILE]`",
	Run: func(cmd *cobra.Command, args []string) {
		e := huffmyfile.Encoder{}
		e.EncodeToDefaultOutputFile(args[0])
	},
}

// Function to return huff command for testing
func NewHuffCmd(testFileName string) *cobra.Command {
	return &cobra.Command{
		Use:   "huff",
		Short: "Compresses .txt files into .huff files. Usage: `huffmyfile huff [FILE]`",
		Run: func(cmd *cobra.Command, args []string) {
			e := huffmyfile.Encoder{}
			e.EncodeToDefaultOutputFile(testFileName)
		},
	}
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
