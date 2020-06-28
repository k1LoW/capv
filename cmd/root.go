/*
Copyright © 2020 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/k1LoW/capv/cap"
	"github.com/spf13/cobra"
)

var (
	pid  int
	path string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "capv",
	Short: "Viewer of Linux capabilitiies.",
	Long:  `Viewer of Linux capabilitiies.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case pid > 0 && path != "":
			printFatalln(cmd, errors.New("not implemented"))
		case pid > 0 && path == "":
			p := cap.NewProc(pid)
			c, err := cap.NewProcCaps(p)
			if err != nil {
				printFatalln(cmd, err)
			}
			if err := c.Pretty(os.Stdout); err != nil {
				printFatalln(cmd, err)
			}
		case pid == 0 && path != "":
			c, err := cap.NewFileCaps(path)
			if err != nil {
				printFatalln(cmd, err)
			}
			if err := c.Pretty(os.Stdout); err != nil {
				printFatalln(cmd, err)
			}
		}
	},
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		printFatalln(rootCmd, err)
	}
}

// https://github.com/spf13/cobra/pull/894
func printErrln(c *cobra.Command, i ...interface{}) {
	c.PrintErr(fmt.Sprintln(i...))
}

func printErrf(c *cobra.Command, format string, i ...interface{}) {
	c.PrintErr(fmt.Sprintf(format, i...))
}

func printFatalln(c *cobra.Command, i ...interface{}) {
	printErrln(c, i...)
	os.Exit(1)
}

func printFatalf(c *cobra.Command, format string, i ...interface{}) {
	printErrf(c, format, i...)
	os.Exit(1)
}

func init() {
	rootCmd.Flags().IntVarP(&pid, "pid", "p", 0, "PID of process")
	rootCmd.Flags().StringVarP(&path, "file", "f", "", "file path")
}
