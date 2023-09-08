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

// unhuffCmd represents the unhuff command
var unhuffCmd = &cobra.Command{
	Use:   "unhuff",
	Short: "Decompresses .huff files into .txt files. Usage: `huffmyfile unhuff [FILE]`",
	Run: func(cmd *cobra.Command, args []string) {
		e := huffmyfile.Encoder{}
		e.DecodeToDefaultOutputFile(args[0])
	},
}

// Wrapper function to return an unhuff command for testing
func NewUnhuffCmd(CompressedTestFileName string) *cobra.Command {
	return &cobra.Command{
		Use:   "unhuff",
		Short: "Decompresses .huff files into .txt files. Usage: `huffmyfile unhuff [FILE]`",
		Run: func(cmd *cobra.Command, args []string) {
			e := huffmyfile.Encoder{}
			e.DecodeToDefaultOutputFile(CompressedTestFileName)
		},
	}
}
func init() {
	rootCmd.AddCommand(unhuffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unhuffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unhuffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
