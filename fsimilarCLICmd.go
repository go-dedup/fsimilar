////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/mkideal/cli"
	"github.com/spakin/awk"
	"gopkg.in/pipe.v2"
	//clix "github.com/mkideal/cli/ext"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The FileT type defines all the necessary info for a file.
type FileT struct {
	Dir  string
	Name string
	Ext  string
	Size int
	Hash uint64
	Dist uint8
}

//==========================================================================
// Main dispatcher

func fsimilar(ctx *cli.Context) error {
	ctx.JSON(ctx.RootArgv())
	ctx.JSON(ctx.Argv())
	fmt.Println()
	rootArgv = ctx.RootArgv().(*rootT)

	Opts.Distance, Opts.SizeGiven, Opts.Template, Opts.Verbose =
		rootArgv.Distance, rootArgv.SizeGiven, rootArgv.Template,
		rootArgv.Verbose.Value()

	return fSimilar(rootArgv.Filei)
}

func fSimilar(cin io.Reader) error {
	psrc, pdst := pipe.Script(
		pipe.Read(cin),
	), pipe.Script(
		pipe.Write(os.Stdout),
	)

	rand.Seed(time.Now().UTC().UnixNano())
	//tmpfile := fmt.Sprintf("%s.%d", file, 99999999-rand.Int31n(90000000))

	p := pipe.Line(
		psrc,
		scripttAwk(),
		pdst,
	)
	err := pipe.Run(p)
	abortOn("Pipe.Run", err)
	println("\nWrapping up.\n")

	return nil
}

func scripttAwk() pipe.Pipe {
	return pipe.TaskFunc(func(st *pipe.State) error {
		// == Setup
		s := awk.NewScript()
		s.Output = st.Stdout

		s.AppendStmt(nil, func(s *awk.Script) {
			file := FileT{}

			if Opts.SizeGiven {
				file.Size = s.F(1).Int()
				// $1 = ""
				// X: s.SetF(1, "")
			}
			p, n := filepath.Split(s.F(0).String())
			file.Dir, file.Name, file.Ext = p, Basename(n), filepath.Ext(n)
			fmt.Printf("  d='%s', n='%s', e='%s', s='%d'\n",
				file.Dir, file.Name, file.Ext, file.Size)
		})

		// 1; # i.e., print all
		//s.AppendStmt(nil, nil)

		// == Run it
		return s.Run(st.Stdin)
	})
}
