////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

//==========================================================================
// fsimilar

type rootT struct {
	cli.Helper
	Distance  uint8        `cli:"d,distance" usage:"the hamming distance of hashes within which to deem similar" dft:"3"`
	SizeGiven bool         `cli:"S,size-given" usage:"size of the files are also available from input (as 1st field)"`
	Filei     *clix.Reader `cli:"*i,input" usage:"input from stdin or the given file (mandatory)"`
	Template  string       `cli:"t,template" usage:"template file name" dft:"fsimilar.tmpl"`
	Verbose   cli.Counter  `cli:"v,verbose" usage:"Verbose mode (Multiple -v options increase the verbosity.)"`
}

var root = &cli.Command{
	Name: "fsimilar",
	Desc: "find/file similar\nVersion " + version + " built on " + date,
	Text: "Find similar files" +
		"\n\nUsage:\n  fsimilar [Options] dir\n\nExample:\n  find . \\( -type f -o -type l \\) -printf '%%7s %%p\\n' | fsimilar -S -i\n  mlocate -i soccer | fsimilar -i",
	Argv: func() interface{} { return new(rootT) },
	Fn:   fsimilar,

	NumOption: cli.AtLeast(1),
}

// Template for main starts here
////////////////////////////////////////////////////////////////////////////
// Global variables definitions

//  var (
//          progname  = "fsimilar"
//          version   = "0.1.0"
//          date = "2017-08-13"
//  )

//  var rootArgv *rootT

////////////////////////////////////////////////////////////////////////////
// Function definitions

// Function main
//  func main() {
//  	cli.SetUsageStyle(cli.ManualStyle) // up-down, for left-right, use NormalStyle
//  	//NOTE: You can set any writer implements io.Writer
//  	// default writer is os.Stdout
//  	if err := cli.Root(root,).Run(os.Args[1:]); err != nil {
//  		fmt.Fprintln(os.Stderr, err)
//  	}
//  	fmt.Println("")
//  }

// Template for main dispatcher starts here
//==========================================================================
// Main dispatcher

//  func fsimilar(ctx *cli.Context) error {
//  	ctx.JSON(ctx.RootArgv())
//  	ctx.JSON(ctx.Argv())
//  	fmt.Println()

//  	return nil
//  }

// Template for CLI handling starts here
