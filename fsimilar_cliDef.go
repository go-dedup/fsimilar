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
	SizeGiven bool         `cli:"S,size-given" usage:"size of the files in input as first field"`
	QuerySize bool         `cli:"Q,query-size" usage:"query the file sizes from os"`
	Filei     *clix.Reader `cli:"*i,input" usage:"input from stdin or the given file (mandatory)"`
	Phonetic  bool         `cli:"p,phonetic" usage:"use phonetic as words for further error tolerant"`
	Final     bool         `cli:"F,final" usage:"produce final output, the recommendations"`
	Ext       string       `cli:"e,ext" usage:"extension to override all files' to (for ffcvt)"`
	CfgPath   string       `cli:"c,cp" usage:"config path, path that hold all template files" dft:"$FSIM_CP"`
	Verbose   cli.Counter  `cli:"v,verbose" usage:"verbose mode (multiple -v increase the verbosity)"`
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
// Constant and data type/structure definitions

// The OptsT type defines all the configurable options from cli.
//  type OptsT struct {
//  	SizeGiven	bool
//  	QuerySize	bool
//  	Filei	*clix.Reader
//  	Phonetic	bool
//  	Final	bool
//  	Ext	string
//  	CfgPath	string
//  	Verbose	cli.Counter
//  	Verbose int
//  }

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

//  var (
//          progname  = "fsimilar"
//          version   = "0.1.0"
//          date = "2017-09-14"

//  	rootArgv *rootT
//  	// Opts store all the configurable options
//  	Opts OptsT
//  )

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
//  	Opts.SizeGiven, Opts.QuerySize, Opts.Filei, Opts.Phonetic, Opts.Final, Opts.Ext, Opts.CfgPath, Opts.Verbose, Opts.Verbose =
//  		rootArgv.SizeGiven, rootArgv.QuerySize, rootArgv.Filei, rootArgv.Phonetic, rootArgv.Final, rootArgv.Ext, rootArgv.CfgPath, rootArgv.Verbose, rootArgv.Verbose.Value()
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
//  	Opts.SizeGiven, Opts.QuerySize, Opts.Filei, Opts.Phonetic, Opts.Final, Opts.Ext, Opts.CfgPath, Opts.Verbose, Opts.Verbose =
//  		rootArgv.SizeGiven, rootArgv.QuerySize, rootArgv.Filei, rootArgv.Phonetic, rootArgv.Final, rootArgv.Ext, rootArgv.CfgPath, rootArgv.Verbose, rootArgv.Verbose.Value()
//  	return nil
//  }

type vecT struct {
	Threshold float64 `cli:"t,thr" usage:"the threshold above which to deem similar (0.8 = 80%%)" dft:"0.86"`
}

var vecDef = &cli.Command{
	Name: "vec",
	Desc: "Use Vector Space for similarity check",
	Text: "Usage:\n  { mlocate -i soccer; mlocate -i football; } | fsimilar sim -i | fsimilar vec -i -S -Q -F",
	Argv: func() interface{} { return new(vecT) },
	Fn:   vecCLI,

	NumOption: cli.AtLeast(1),
}
