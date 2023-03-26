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
	"strings"

	"github.com/mewkiz/pkg/jsonutil"
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
	// parse command line arguments.
	var listfilePath string
	flag.StringVar(&listfilePath, "listfile", "", "listfile path (JSON format)")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	listfile, err := parseListfile(listfilePath)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	for _, pakPath := range flag.Args() {
		if err := dumpPakArchive(pakPath, rootDumpDir, listfile); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

const (
	// rootDumpDir specifies the top-level output directory.
	rootDumpDir = "_dump_"
	// keepSubarchive specifies whether to keep PAK (sub)archives after
	// extracting files.
	//
	// Note: the top level PAK archive is always kept.
	keepSubarchive = false
)

// dumpPakArchive dumps the given PAK archive to the specified output directory.
func dumpPakArchive(pakPath, dumpDir string, listfile map[string]string) error {
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
	reverseListfile := make(map[string]string)
	for key, val := range listfile {
		reverseListfile[val] = key
	}
	// output PAK subarchives (and files).
	var subarchivePaths []string
	for i, fileContents := range filesContents {
		name := "archive"
		ext := "pak"
		if !isArchive(fileContents) {
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
		if len(fileContents) == 0 {
			//dbg.Println("skip empty file %q", dstPath)
			continue
		}
		if len(listfile) > 0 {
			if newPathName, ok := replaceName(dstPath, listfile); ok {
				dstPath = newPathName
			} else {
				warn.Printf("file name not set for %q in listfile", dstPath)
			}
		}
		dbg.Printf("creating %q", dstPath)
		if err := ioutil.WriteFile(dstPath, fileContents, 0o644); err != nil {
			return errors.WithStack(err)
		}
		if filepath.Ext(dstPath) == ".pak" {
			subarchivePaths = append(subarchivePaths, dstPath)
		}
	}
	// dump subarchives.
	//dbg.Printf("--- [ dumping subarchives of %q ] ---", pakPath)
	for _, subarchivePath := range subarchivePaths {
		if err := dumpPakArchive(subarchivePath, dstDir, listfile); err != nil {
			return errors.WithStack(err)
		}
		// Only remove subarchive if present in listfile. If not present, it's
		// probably a ZEL file that's not yet successfully decoded by zel_dump or
		// a subdirectory that has not yet been named.
		_, inListfile := reverseListfile[stripRootDumpDir(subarchivePath)]
		if !keepSubarchive && inListfile {
			// only keep extracted files of subarchive.
			if err := os.Remove(subarchivePath); err != nil {
				return errors.WithStack(err)
			}
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

// parseListfile parses the given listfile.
//
// Example listfile:
//
//	{
//		"X/archive_0000.pak": "X/core.pak",
//		"X/core/file_0002.bin": "X/core/core.pal",
//		"X/core/file_0003.bin": "X/core/palette.bmp"
//	}
func parseListfile(listfilePath string) (map[string]string, error) {
	listfile := make(map[string]string)
	if len(listfilePath) == 0 {
		return listfile, nil
	}
	if err := jsonutil.ParseFile(listfilePath, &listfile); err != nil {
		return nil, errors.WithStack(err)
	}
	return listfile, nil
}

// replaceName replaces the given path by a corresponding new path if a
// replacement was specified in the given listfile.
func replaceName(path string, listfile map[string]string) (string, bool) {
	if newPath, ok := listfile[stripRootDumpDir(path)]; ok {
		return filepath.Join(rootDumpDir, newPath), true
	}
	return path, false
}

// stripRootDumpDir strips the root dump directory ("_dump_") from the prefix of
// the given path.
func stripRootDumpDir(path string) string {
	parts := strings.Split(path, string(filepath.Separator))
	if parts[0] != rootDumpDir {
		panic(fmt.Errorf("invalid path root; expected root dump dir %q, got %q", rootDumpDir, parts[0]))
	}
	return filepath.Join(parts[1:]...)
}
