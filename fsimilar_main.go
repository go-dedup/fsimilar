////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

//go:generate sh -v fsimilar_cliGen.sh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-dedup/text"
	"github.com/labstack/gommon/color"
	"github.com/mkideal/cli"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The OptsT type defines all the configurable options from cli.
type OptsT struct {
	SizeGiven bool
	QuerySize bool
	Phonetic  bool
	Final     bool
	Template  string
	ExecPath  string
	CfgPath   string
	Verbose   int
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	progname = "fsimilar"
	version  = "0.1.0"
	date     = "2017-09-02"
)

var (
	rootArgv *rootT
	// Opts store all the configurable options for jsonfiddle.
	Opts OptsT

	doc2words = text.GetWordsFactory(text.Decorators(
		text.SplitCamelCase,
		text.ToLower,
		text.RemovePunctuation,
		text.Compact,
		text.Trim,
	))
)

////////////////////////////////////////////////////////////////////////////
// Function definitions

//==========================================================================
// Main dispatcher

func fsimilar(ctx *cli.Context) error {
	ctx.JSON(ctx.RootArgv())
	ctx.JSON(ctx.Argv())
	fmt.Println()
	return nil
}

//==========================================================================
// Main

func main() {
	Opts.ExecPath = filepath.Dir(os.Args[0])
	// cli.SetUsageStyle(cli.ManualStyle)
	if err := cli.Root(root,
		cli.Tree(simDef),
		cli.Tree(vecDef)).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("")
}

//==========================================================================
// support functions

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Basename returns the file name without extension.
func Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

// IsExist checks if the given file exist
func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func warning(m string) {
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", progname, color.Yellow("Warning"), m)
}

func warnOn(errCase string, e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s, %s: %v\n",
			color.White(progname), color.Yellow("Error"), errCase, e)
	}
}

// abortOn will quit on anticipated errors gracefully without stack trace
func abortOn(errCase string, e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s, %s: %v\n",
			color.White(progname), color.Red("Error"), errCase, e)
		os.Exit(1)
	}
}

// verbose will print info to stderr according to the verbose level setting
func verbose(levelSet int, format string, args ...interface{}) {
	if Opts.Verbose >= levelSet {
		fmt.Fprintf(os.Stderr, "["+progname+"] "+format+"\n", args...)
	}
}
