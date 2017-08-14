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
	"time"

	"github.com/mkideal/cli"
	"github.com/spakin/awk"
	"gopkg.in/pipe.v2"
	//clix "github.com/mkideal/cli/ext"
)

//==========================================================================
// Main dispatcher

func fsimilar(ctx *cli.Context) error {
	ctx.JSON(ctx.RootArgv())
	ctx.JSON(ctx.Argv())
	fmt.Println()
	rootArgv = ctx.RootArgv().(*rootT)

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

		// 1; # i.e., print all
		s.AppendStmt(nil, nil)

		// == Run it
		return s.Run(st.Stdin)
	})
}
