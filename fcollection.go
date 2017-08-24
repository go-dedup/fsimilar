////////////////////////////////////////////////////////////////////////////
// Program: fcollection
// Purpose: file collections
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

//==========================================================================

// The FileT type defines all the necessary info for a file.
type FileT struct {
	Name string
	Ext  string
	Size int
	Dir  string
	Hash uint64
	Dist uint8
	Vstd bool // Visited

	SizeRef int
}

func (f FileT) Similarity() int {
	return int((1-float32(Abs(f.Size-f.SizeRef))/float32(f.SizeRef))*
		(1-float32(f.Dist)/float32(64))*100 + 0.5)
}

type Files []FileT

func (a Files) Len() int      { return len(a) }
func (a Files) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Files) Less(i, j int) bool {
	if a[i].Name < a[j].Name {
		return true
	}
	if a[i].Name > a[j].Name {
		return false
	}
	if a[i].Dist < a[j].Dist {
		return true
	}
	if a[i].Dist > a[j].Dist {
		return false
	}
	return a[i].Size < a[j].Size
}

//==========================================================================

// The HVisited tells whether the hash has been visited.
type HVisited map[uint64]bool

//==========================================================================

// The FCollection type defines the collections between hashes and files.
type FCollection map[uint64]Files

func NewFCollection() FCollection {
	return make(FCollection)
}

func (fc FCollection) Add(key uint64, value FileT) {
	fc[key] = append(fc[key], value)
}

func (fc FCollection) Get(key uint64) (Files, bool) {
	val, ok := fc[key]
	return val, ok
}

func (fc FCollection) Len() int {
	return len(fc)
}

func (fc FCollection) LenOf(key uint64) int16 {
	return int16(len(fc[key]))
}

//==========================================================================

// The FCItem type can be used to identify an item in the collections.
type FCItem struct {
	Hash  uint64
	Index int16
}

//==========================================================================

// The FAll type defines the lookup from file to FCItem.
type FAll map[string]FCItem
