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
	SizeGiven bool         `cli:"S,size-given" usage:"size of the files available from input as 1st field"`
	QuerySize bool         `cli:"Q,query-size" usage:"file size not available so query it from filesystem"`
	Filei     *clix.Reader `cli:"*i,input" usage:"input from stdin or the given file (mandatory)"`
	Phonetic  bool         `cli:"p,phonetic" usage:"use phonetic to group sound-similar words for further error tolerant"`
	Verbose   cli.Counter  `cli:"v,verbose" usage:"verbose mode (multiple -v options increase the verbosity)"`
}

var root = &cli.Command{
	Name:   "fsimilar",
	Desc:   "find/file similar\nVersion " + version + " built on " + date,
	Text:   "Find similar files",
	Global: true,
	Argv:   func() interface{} { return new(rootT) },
	Fn:     fsimilar,

	NumOption: cli.AtLeast(1),
}

// Template for main starts here
////////////////////////////////////////////////////////////////////////////
// Global variables definitions

//  var (
//          progname  = "fsimilar"
//          version   = "0.1.0"
//          date = "2017-08-31"
//  )

//  var rootArgv *rootT

////////////////////////////////////////////////////////////////////////////
// Function definitions

// Function main
//  func main() {
//  	cli.SetUsageStyle(cli.ManualStyle) // up-down, for left-right, use NormalStyle
//  	//NOTE: You can set any writer implements io.Writer
//  	// default writer is os.Stdout
//  	if err := cli.Root(root,
//  		cli.Tree(simDef),
//  		cli.Tree(vecDef)).Run(os.Args[1:]); err != nil {
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

////////////////////////////////////////////////////////////////////////////
// sim

//  func simCLI(ctx *cli.Context) error {
//  	rootArgv = ctx.RootArgv().(*rootT)
//  	argv := ctx.Argv().(*simT)
//  	fmt.Printf("[sim]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
//  	return nil
//  }

type simT struct {
	Distance uint8 `cli:"d,dist" usage:"the hamming distance of hashes within which to deem similar" dft:"3"`
}

var simDef = &cli.Command{
	Name: "sim",
	Desc: "Filter the input using simhash similarity check",
	Text: "Usage:\n  mlocate -i soccer | fsimilar sim -i",
	Argv: func() interface{} { return new(simT) },
	Fn:   simCLI,

	NumOption: cli.AtLeast(1),
}

////////////////////////////////////////////////////////////////////////////
// vec

//  func vecCLI(ctx *cli.Context) error {
//  	rootArgv = ctx.RootArgv().(*rootT)
//  	argv := ctx.Argv().(*vecT)
//  	fmt.Printf("[vec]:\n  %+v\n  %+v\n  %v\n", rootArgv, argv, ctx.Args())
//  	return nil
//  }

type vecT struct {
	Threshold float64 `cli:"t,threshold" usage:"the threshold above which to deem similar" dft:"0.6"`
}

var vecDef = &cli.Command{
	Name: "vec",
	Desc: "Use Vector Space for similarity check",
	Text: "Usage:\n  mlocate -i soccer | fsimilar sim -i | fsimilar vec -i -Q",
	Argv: func() interface{} { return new(vecT) },
	Fn:   vecCLI,

	NumOption: cli.AtLeast(1),
}
