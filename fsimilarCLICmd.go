////////////////////////////////////////////////////////////////////////////
// Program: fsimilar
// Purpose: find/file similar
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/go-dedup/simhash"
	"github.com/go-dedup/simhash/sho"

	"github.com/mkideal/cli"
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
	Vstd bool // Visited
}

type Files []FileT

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
	rand.Seed(time.Now().UTC().UnixNano())
	//tmpfile := fmt.Sprintf("%s.%d", file, 99999999-rand.Int31n(90000000))

	oracle := sho.NewOracle()
	sh := simhash.NewSimhash()
	r := Opts.Distance
	fAll := make(FAll)     // the all file to FCItem map
	fc := NewFCollection() // the FCollection that holds everything

	// read input line by line
	scanner := bufio.NewScanner(cin)
	for scanner.Scan() {
		file := FileT{}
		fn := scanner.Text()

		// == Gather file info
		if Opts.SizeGiven {
			_, err := fmt.Sscan(fn, &file.Size)
			abortOn("Parsing file size", err)
			il := regexp.MustCompile(`^ *\d+\s+(.*)$`).FindStringSubmatchIndex(fn)
			//fmt.Println(il)
			fn = fn[il[2]:]
		} else {
			s, err := os.Stat(fn)
			warnOn("Get file size", err)
			file.Size = int(s.Size())
		}
		p, n := filepath.Split(fn)
		file.Dir, file.Name, file.Ext = p, Basename(n), filepath.Ext(n)
		verbose(1, " n='%s', e='%s', s='%d', d='%s'",
			file.Name, file.Ext, file.Size, file.Dir)

		hash := sh.GetSimhash(sh.NewWordFeatureSet([]byte(file.Name)))
		fi := FCItem{Hash: hash, Index: fc.LenOf(hash)}
		fAll[fn] = fi
		fc.Add(hash, file)

		// == Build similarity knowledge
		if h, d, seen := oracle.Find(hash, r); seen == true {
			verbose(1, "=: Simhash of %x ignored for %x (%d).", hash, h, d)
		} else {
			oracle.See(hash)
			verbose(1, "+: Simhash of %x added.", hash)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// process all, the sorted fAll map
	visited := make(HVisited)
	var keys []string
	for k := range fAll {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// for each sorted item according to file path & name
	for _, k := range keys {
		// get next file item from FAll file to FCItem map
		fi := fAll[k]
		// and skip if the hash has already been visited
		if visited[fi.Hash] {
			continue
		}
		visited[fi.Hash] = true
		// also skip if no similar items at this hash
		if fc.LenOf(fi.Hash) <= 1 {
			continue
		}

		// similarity exist, start digging
		files, ok := fc.Get(fi.Hash)
		if !ok {
			abortOn("Internal error", errors.New("fc integrity checking"))
		}
		for ii, _ := range files {
			files[ii].Dist = 0
		}
		for _, nigh := range oracle.Search(fi.Hash, r) {
			visited[nigh.H] = true
			// files more
			fm, ok := fc.Get(nigh.H)
			if ok {
				for ii, _ := range fm {
					fm[ii].Dist = nigh.D
				}
				files = append(files, fm...)
			}
		}

		// One group of similar items found, output
		verbose(1, "## Similar items\n %v.", files)
	}

	return nil
}
