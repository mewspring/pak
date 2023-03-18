// PAK file format
//
//    archive_offsets [...]uint32
//    archive_data    []byte

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mewkiz/pkg/pathutil"
	"github.com/mewkiz/pkg/term"
	"github.com/mewspring/pak/archive/pak"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "pak_dump:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("pak_dump:")+" ", 0)
	// warn is a logger with the "pak_dump:" prefix which logs warning messages
	// to standard error.
	warn = log.New(os.Stderr, term.RedBold("pak_dump:")+" ", log.Lshortfile)
)

func usage() {
	const usage = "Usage: pak_dump [OPTIONS]... FILE.pak..."
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	for _, pakPath := range flag.Args() {
		if err := dumpPakArchive(pakPath, rootDumpDir); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

// rootDumpDir specifies the top-level output directory.
const rootDumpDir = "_dump_"

// dumpPakArchive dumps the given PAK archive to the specified output directory.
func dumpPakArchive(pakPath, dumpDir string) error {
	// parse PAK archive.
	dbg.Printf("extracting %q", pakPath)
	filesContents, err := pak.Extract(pakPath)
	if err != nil {
		return errors.WithStack(err)
	}
	// create output directory.
	pakName := pathutil.FileName(pakPath)
	pakNameWithoutExt := pathutil.TrimExt(pakName)
	dstDir := filepath.Join(dumpDir, pakNameWithoutExt)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return errors.WithStack(err)
	}
	// output PAK subarchives (and files).
	var subarchivePaths []string
	for i, fileContents := range filesContents {
		name := "archive"
		ext := "pak"
		isArch := isArchive(fileContents)
		if !isArch {
			if isSound(fileContents) {
				name = "sound"
				ext = "wav"
			} else {
				name = "file"
				ext = "bin"
			}
		}
		dstName := fmt.Sprintf("%s_%04d.%s", name, i, ext)
		dstPath := filepath.Join(dstDir, dstName)
		dbg.Printf("creating %q", dstPath)
		if err := ioutil.WriteFile(dstPath, fileContents, 0o644); err != nil {
			return errors.WithStack(err)
		}
		if isArch {
			subarchivePaths = append(subarchivePaths, dstPath)
		}
	}
	// dump subarchives.
	//dbg.Printf("--- [ dumping subarchives of %q ] ---", pakPath)
	for _, subarchivePath := range subarchivePaths {
		if err := dumpPakArchive(subarchivePath, dstDir); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// isArchive reports whether the given contents is a PAK archive.
func isArchive(buf []byte) bool {
	_, err := pak.ParsePAKHeader(buf)
	return err == nil
}

// isSound reports whether the given contents is a WAV sound file.
func isSound(buf []byte) bool {
	if len(buf) < 4 {
		return false
	}
	return string(buf[0:4]) == "RIFF"
}
