// PAK file format
//
//    archive_offsets [...]uint32
//    archive_data    []byte

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/pathutil"
	"github.com/mewkiz/pkg/term"
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

func main() {
	flag.Parse()
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
	// read PAK file contents.
	buf, err := ioutil.ReadFile(pakPath)
	if err != nil {
		return errors.WithStack(err)
	}
	if !isArchive(buf) {
		//dbg.Printf("not an archive: %q", pakPath)
		return nil
	}
	dbg.Printf("extracting %q", pakPath)
	// parse PAK header.
	pakHdrSize := int(binary.LittleEndian.Uint32(buf[0:4]))
	r := bytes.NewReader(buf)
	pakHdrReader := io.NewSectionReader(r, 0, int64(pakHdrSize))
	narchives := pakHdrSize / 4
	archiveOffsets := make([]uint32, narchives)
	if err := binary.Read(pakHdrReader, binary.LittleEndian, &archiveOffsets); err != nil {
		return errors.WithStack(err)
	}
	if len(buf) != int(archiveOffsets[len(archiveOffsets)-1]) {
		pretty.Println("archiveOffsets:", archiveOffsets)
		panic(fmt.Errorf("mismatch between archiveOffsets[%d]=%d and len(buf)=%d", len(archiveOffsets)-1, archiveOffsets[len(archiveOffsets)-1], len(buf)))
	}
	// create output directory.
	pakName := pathutil.FileName(pakPath)
	pakNameWithoutExt := pathutil.TrimExt(pakName)
	dstDir := filepath.Join(dumpDir, pakNameWithoutExt)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return errors.WithStack(err)
	}
	// output PAK subarchives (and files).
	var archivePaths []string
	for i := 0; i < len(archiveOffsets)-1; i++ {
		archiveStartOffset := archiveOffsets[i]
		archiveEndOffset := archiveOffsets[i+1]
		archiveContents := buf[archiveStartOffset:archiveEndOffset]
		name := "archive"
		if !isArchive(archiveContents) {
			name = "file"
		}
		archiveName := fmt.Sprintf("%s_%04d.bin", name, i)
		archivePath := filepath.Join(dstDir, archiveName)
		dbg.Printf("creating %q", archivePath)
		if err := ioutil.WriteFile(archivePath, archiveContents, 0o644); err != nil {
			return errors.WithStack(err)
		}
		archivePaths = append(archivePaths, archivePath)
	}
	// dump subarchives.
	//dbg.Printf("--- [ dumping subarchives of %q ] ---", pakPath)
	for _, archivePath := range archivePaths {
		if err := dumpPakArchive(archivePath, dstDir); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// isArchive reports whether the given contents is a PAK archive.
func isArchive(buf []byte) bool {
	if len(buf) < 4 {
		return false
	}
	// parse PAK header.
	pakHdrSize := int(binary.LittleEndian.Uint32(buf[0:4]))
	// the minimum valid PAK archive header is 8 bytes as a start and end offset
	// is required for each file. a PAK archive containing a single empty file
	// would have the PAK header `00 00 00 00  08 00 00 00`.
	if pakHdrSize < 8 || pakHdrSize > len(buf) {
		return false
	}
	r := bytes.NewReader(buf)
	pakHdrReader := io.NewSectionReader(r, 0, int64(pakHdrSize))
	narchives := pakHdrSize / 4
	archiveOffsets := make([]uint32, narchives)
	if err := binary.Read(pakHdrReader, binary.LittleEndian, &archiveOffsets); err != nil {
		panic(fmt.Errorf("unable to read PAK header; %+v", err))
	}
	if len(buf) != int(archiveOffsets[len(archiveOffsets)-1]) {
		return false
	}
	return true
}
