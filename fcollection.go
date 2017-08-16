////////////////////////////////////////////////////////////////////////////
// Program: fcollection
// Purpose: file collections
// Authors: Tong Sun (c) 2017, All rights reserved
////////////////////////////////////////////////////////////////////////////

package main

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

//==========================================================================

// The HVisited tells whether the hash has been visited.
type HVisited map[uint64]bool

//==========================================================================

// The FCollection type defines the collections between hashes and files.
type FCollection map[uint64]Files

func NewFCollection() FCollection {
	return make(FCollection)
}

func (fc FCollection) Set(key uint64, value FileT) {
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

//==========================================================================

// The FClose type defines a pair of hashes and files that are close.
type FClose struct {
	Item FCItem // The duplicated
	Hash uint64 // that close to this hash
	Dist uint8  // with this distance
}

//==========================================================================

// The FCs type is the slice of FClose.
type FCs []FClose

func (fcs FCs) SetClose(item FCItem, hash uint64, dist uint8) {
	fcs = append(fcs, FClose{item, hash, dist})
}
