package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "zel_patch:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("zel_patch:")+" ", 0)
	// warn is a logger with the "zel_patch:" prefix which logs warning messages
	// to standard error.
	warn = log.New(os.Stderr, term.RedBold("zel_patch:")+" ", log.Lshortfile)
)

func usage() {
	const usage = "Usage: zel_patch [OPTIONS]..."
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	// patch files.
	for _, file := range files {
		if err := patch(file); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

// patch patches the given file.
func patch(file File) error {
	buf, err := ioutil.ReadFile(file.path)
	if err != nil {
		return errors.WithStack(err)
	}
	rawHash := sha1.Sum(buf)
	hash := fmt.Sprintf("%040x", rawHash[:])
	switch hash {
	case file.hashBefore:
		// nothing to do; expected case before patching.
	case file.hashAfter:
		// already patched, early return.
		dbg.Printf("file %q already patched", file.path)
		return nil
	default:
		return errors.Errorf("unable to patch file %q with unexpected contents; expected hash %s, got %s", file.path, file.hashBefore, hash)
	}
	var data []byte
	off := 0
	for _, replace := range file.replaces {
		data = append(data, buf[off:replace.pos]...)
		off = replace.pos
		data = append(data, replace.after...)
		off += len(replace.before)
	}
	data = append(data, buf[off:]...)
	dbg.Printf("patching %q", file.path)
	if err := ioutil.WriteFile(file.path, data, 0o644); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// File specifies a file to patch.
type File struct {
	// Path to file.
	path string
	// SHA1 hash of file contents before patch.
	hashBefore string
	// SHA1 hash of file contents after patch.
	hashAfter string
	// Replacements.
	replaces []Replace
}

// Replace specifies the before and after for a given position of the file.
type Replace struct {
	// File offset.
	pos int
	// Contents before patch.
	before []byte
	// Contents after patch.
	after []byte
}

// files specifies the files to patch.
var files = []File{
	// patch to fix frame 221 in X/tilesets/tileset_4_buildings.zel.
	{
		path:       "_dump_/X/tilesets/tileset_4_buildings.zel",
		hashBefore: "776a9f27489da08bcd85b654eaf0474f90994449",
		hashAfter:  "6c74668a0d168c49b3f33b08b1c93dd8ab072fe7",
		replaces: []Replace{
			// remove two extra bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x32D817,
				before: []byte{0x1B, 0x1B}, // NOTE: which two pixels to remove is unknown.
				after:  []byte{},
			},
			// add two missing bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x32DB71,
				before: []byte{},
				after:  []byte{0x1B, 0x1B}, // NOTE: which two pixels to add is unknown.
			},
			// remove two extra bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x32E73C,
				before: []byte{0x1C, 0x0E}, // NOTE: which two pixels to remove is unknown.
				after:  []byte{},
			},
			// add two missing bytes of cmd directive.
			{
				pos:    0x32EB78,
				before: []byte{},
				// NOTE: adding 0x002E cmd to match format of previous line
				//
				//	prev_line
				//		cmd=0x1007 // regular pixels (npixels=7)
				//		cmd=0x002C // transparent pixels (xSkip=44)
				//		cmd=0x900D // regular pixels (npixels=13)
				//
				//		7+44+13 = 64 (frame width)
				//
				// this_line (after patch)
				//		cmd=0x1007 // regular pixels (npixels=7)
				//		cmd=0x002E // regular pixels (npixels=46) <<--- added
				//		cmd=0x900B // regular pixels (npixels=11)
				//
				//		7+46+11 = 64 (frame width)
				after: []byte{0x2E, 0x00},
			},
			// remove two extra bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x32F847,
				before: []byte{0x3E, 0x3E}, // NOTE: which two pixels to remove is unknown.
				after:  []byte{},
			},
			// add two missing bytes of cmd directive.
			{
				pos:    0x32FB97,
				before: []byte{},
				// NOTE: adding 0x8019 cmd to match format of next line
				//
				// this_line (after patch)
				//		cmd=0x000C // transparent pixels (xSkip=12)
				//		cmd=0x101B // regular pixels (npixels=27)
				//		cmd=0x8019 // transparent pixels (xSkip=25) <--- added
				//
				//		12+27+25 = 64 (frame width)
				//
				// next_line
				//		cmd=0x000D // transparent pixels (xSkip=13)
				//		cmd=0x101A // regular pixels (npixels=26)
				//		cmd=0x8019 // transparent pixels (xSkip=25)
				//
				//		13+26+25 = 64 (frame width)
				after: []byte{0x19, 0x80},
			},
		},
	},
	// patch to fix frame 327 and 328 in X/tilesets/tileset_8_buildings.zel.
	{
		path:       "_dump_/X/tilesets/tileset_8_buildings.zel",
		hashBefore: "5b34a4b0f4722b50e461aeba963e37ac85460112",
		hashAfter:  "485d44e59ce719269c52c910d7bd06b824c4b82c",
		replaces: []Replace{
			// --- [ frame 327 ] ---
			//
			// remove two extra bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x30377A,
				before: []byte{0xF9, 0x2D}, // NOTE: which two pixels to remove is unknown.
				after:  []byte{},
			},
			// add two missing bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x303AEF,
				before: []byte{},
				after:  []byte{0xF9, 0x2D}, // NOTE: which two pixels to add is unknown.
			},
			// --- [ frame 328 ] ---
			//
			// remove two extra bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x3047AC,
				before: []byte{0x2F, 0x99}, // NOTE: which two pixels to remove is unknown.
				after:  []byte{},
			},
			// add two missing bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x304B06,
				before: []byte{},
				after:  []byte{0x2F, 0x99}, // NOTE: which two pixels to add is unknown.
			},
		},
	},
	// patch to fix frame 113 in X/tilesets/tileset_14_buildings.zel.
	{
		path:       "_dump_/X/tilesets/tileset_14_buildings.zel",
		hashBefore: "4206e376b52c60990ddfb80c64d0adc5f7d34b66",
		hashAfter:  "4a9a5ca262f98cbef167dd1d053caf4c8007cca0",
		replaces: []Replace{
			// remove two extra bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x1579C8,
				before: []byte{0x0B, 0x0B}, // NOTE: which two pixels to remove is unknown.
				after:  []byte{},
			},
			// add two missing bytes of pixel line (otherwise cmd offset gets skewed).
			{
				pos:    0x157B54,
				before: []byte{},
				after:  []byte{0x0B, 0x0B}, // NOTE: which two pixels to add is unknown.
			},
		},
	},
}
