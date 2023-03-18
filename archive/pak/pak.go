// Package pak provides access to PAK archives.
package pak

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Extract extracts the given PAK archive, returning the contents of the
// top-level files (and subarchives) contained within the archive.
func Extract(pakPath string) ([][]byte, error) {
	// read PAK file contents.
	buf, err := ioutil.ReadFile(pakPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// parse PAK header.
	archiveOffsets, err := ParsePAKHeader(buf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// return contents of files contained within archive.
	narchives := len(archiveOffsets) - 1
	var filesContents [][]byte
	for i := 0; i < narchives; i++ {
		startOffset := archiveOffsets[i]
		endOffset := archiveOffsets[i+1]
		fileContents := buf[startOffset:endOffset:endOffset]
		filesContents = append(filesContents, fileContents)
	}
	return filesContents, nil
}

// ParsePAKHeader parses the PAK header of the given PAK file contents, and
// returns the archive offsets.
func ParsePAKHeader(buf []byte) ([]uint32, error) {
	if len(buf) < 4 {
		return nil, errors.Errorf("too short PAK header; expected >= 4, got %d", len(buf))
	}
	pakHdrSize := int(binary.LittleEndian.Uint32(buf[0:4]))
	// the minimum valid PAK archive header is 8 bytes as a start and end offset
	// is required for each file. a PAK archive containing a single empty file
	// would have the PAK header `00 00 00 00  08 00 00 00`.
	if pakHdrSize < 8 || pakHdrSize > len(buf) {
		return nil, errors.Errorf("invalid PAK header size; expected >= 8 and <= len(buf)=%d, got %d", len(buf), pakHdrSize)
	}
	r := bytes.NewReader(buf)
	pakHdrReader := io.NewSectionReader(r, 0, int64(pakHdrSize))
	archiveOffsetsLen := pakHdrSize / 4
	archiveOffsets := make([]uint32, archiveOffsetsLen)
	if err := binary.Read(pakHdrReader, binary.LittleEndian, &archiveOffsets); err != nil {
		return nil, errors.WithStack(err)
	}
	narchives := len(archiveOffsets) - 1
	if len(buf) != int(archiveOffsets[narchives]) {
		return nil, errors.Errorf("mismatch between archiveOffsets[%d]=%d and len(buf)=%d", narchives, archiveOffsets[narchives], len(buf))
	}
	return archiveOffsets, nil
}
