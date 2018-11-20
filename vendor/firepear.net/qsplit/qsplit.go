/*
Package qsplit (short for "quoted split") performs a Unix shell style
split-on-whitespace of its input. Its functions return the
non-whitespace "chunks" contained in their input, treating text within
balanced quotes as a single chunk.

Whitespace, according to qsplit, is the ASCII space and horizontal tab
characters. qsplit is aware of several quote character pairs:

    ASCII::     '', "", ``
    Guillemets: ‹›, «»
    Japanese:   「」,『』

These are the rules used to delineate chunks:

    * Quotes begin only at a word boundary
    * Quotes extend to the first closing quotation mark which matches
      the opening quote, which may or may not be at a word boundary.
    * Quotes do not nest

*/
package qsplit // import "firepear.net/qsplit"

// Copyright (c) 2014-2016 Shawn Boyette <shawn@firepear.net>. All
// rights reserved.  Use of this source code is governed by a
// BSD-style license that can be found in the LICENSE file.

import (
	"bytes"
	"unicode/utf8"
)

var (
	// Version is the current version
	Version = "2.2.2"

	// the quotation marks we know about
	quotes = map[rune]rune{
		'\'': '\'', '"': '"', '`': '`',
		'‹': '›', '«': '»',
		'「': '」', '『': '』',
	}
)

// Locations finds and returns the beginning and end points of all
// text chunks in its input.
func Locations(b []byte) [][2]int {
	return realLocations(b, false)
}

// LocationsOnce finds and returns only the beginning and end point of
// the first chunk, and the beginning of the next chunk. If this is
// all you need, LocationsOnce is significantly faster than
// Locations.
//
// If no chunks are found, the first element of the returned array
// will be -1. Similarly, if only one chunk is found, the third
// element will be -1.
func LocationsOnce(b []byte) [3]int {
	s := realLocations(b, true)
	var locs [3]int
	if len(s) == 0 {
		locs[0] = -1
	} else {
		locs[0] = s[0][0]
		locs[1] = s[0][1]
		if len(s) == 2 {
			locs[2] = s[1][0]
		} else {
			locs[2] = -1
		}
	}
	return locs
}

// realLocations does the work for Locations and LocationsOnce
func realLocations(b []byte, once bool) [][2]int {
	var si [][2]int       // slice of tuples of ints
	var inw, inq, ok bool // in-word, in-quote, escape flags; map test var
	var rune, endq rune   // current rune; end-quote for current quote
	var i, idx int        // first index of chunk; byte index of current rune

	// we need to operate at the runes level
	runes := bytes.Runes(b)
	for _, rune = range runes {
		switch {
		case inq:
			// in a quoted chunk, if we're looking at the ending
			// quote, unset inq and append a the tuple for this chunk
			// to the return list.
			if rune == endq {
				inq = false
				si = append(si, [2]int{i, idx})
			}
		case rune == ' ' || rune == '\t':
			// if looking at a space and inw is set, end the present
			// chunk and append a new tuple. else just move on.
			if inw {
				inw = false
				si = append(si, [2]int{i, idx})
			}
		case inw:
			// if in a regular word, do nothing
		default:
			if endq, ok = quotes[rune]; ok {
				// looking at an unescaped quote; set inq and i
				inq = true
				i = idx + utf8.RuneLen(rune)
			} else {
				// looking at the first rune in a word. set inw & i
				inw = true
				i = idx
			}
		}
		// if once-mode is on and we've found 2 chunks, return
		if len(si) == 2 && once {
			return si
		}
		// else update idx and prune
		idx += utf8.RuneLen(rune)
	}
	// append the tuple for the last chunk if we were still in a word
	// or quote
	if inw || inq {
		si = append(si, [2]int{i, idx})
	}
	return si
}

// ToBytes performs a quoted split to a slice of byteslices.
func ToBytes(b []byte) [][]byte {
	var sb [][]byte    // slice of slice of bytes
	cp := Locations(b) // get chunk positions
	for _, pos := range cp {
		sb = append(sb, b[pos[0]:pos[1]])
	}
	return sb
}

// ToStrings performs a quoted split to a slice of strings.
func ToStrings(b []byte) []string {
	var ss []string
	cp := Locations(b) // get chunk positions
	for _, pos := range cp {
		ss = append(ss, string(b[pos[0]:pos[1]]))
	}
	return ss
}

// ToStringBytes performs a quoted split, returning the first chunk as
// a string and the rest as a slice of byteslices.
func ToStringBytes(b []byte) (string, [][]byte) {
	bslices := ToBytes(b)
	return string(bslices[0]), bslices[1:]
}

// Once performs a single quoted split, returning the first chunk
// found in the input byteslice, and the remainder of the byteslice
func Once(b []byte) [][]byte {
	var sb [][]byte    // slice of slice of bytes
	cp := Locations(b) // get chunk positions
	if len(cp) == 1 {
		sb = append(sb, b)
	} else {
		sb = append(sb, b[cp[0][0]:cp[0][1]])
		sb = append(sb, b[cp[1][0]:])
	}
	return sb
}
