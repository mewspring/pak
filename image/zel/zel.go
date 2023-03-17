package zel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "zel:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("zel:")+" ", 0)
	// warn is a logger with the "zel:" prefix which logs warning messages to
	// standard error.
	warn = log.New(os.Stderr, term.RedBold("zel:")+" ", log.Lshortfile)
)

// ZEL frame format
//
//    width  uint16
//    height uint16
//    data   []byte

// DecodeAll decodes the given ZEL image using colours from the provided
// palette, and returns the sequential frames.
func DecodeAll(zelPath string, pal color.Palette) ([]image.Image, error) {
	buf, err := ioutil.ReadFile(zelPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("parsing %q", zelPath)
	// parse ZEL header.
	zelHdrSize := int(binary.LittleEndian.Uint32(buf[0:4]))
	r := bytes.NewReader(buf)
	zelHdrReader := io.NewSectionReader(r, 0, int64(zelHdrSize))
	nframes := zelHdrSize / 4
	frameOffsets := make([]uint32, nframes)
	if err := binary.Read(zelHdrReader, binary.LittleEndian, &frameOffsets); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(buf) != int(frameOffsets[len(frameOffsets)-1]) {
		pretty.Println("frameOffsets:", frameOffsets)
		panic(fmt.Errorf("mismatch between frameOffsets[%d]=%d and len(buf)=%d", len(frameOffsets)-1, frameOffsets[len(frameOffsets)-1], len(buf)))
	}
	// output ZEL frames.
	var imgs []image.Image
	for i := 0; i < len(frameOffsets)-1; i++ {
		frameStartOffset := frameOffsets[i]
		frameEndOffset := frameOffsets[i+1]
		frameContents := buf[frameStartOffset:frameEndOffset]
		img, ok := parseFrame(frameContents, pal)
		if !ok {
			warn.Printf("skipping invalid frame (%d/%d) of %q", i, len(frameOffsets)-1, zelPath)
			continue // skip
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
}

// parseFrame parses the given ZEL frame contents.
func parseFrame(frameContents []byte, pal color.Palette) (image.Image, bool) {
	// parse ZEL frame.
	frameWidth := int(binary.LittleEndian.Uint16(frameContents[0:2]))
	frameHeight := int(binary.LittleEndian.Uint16(frameContents[2:4]))
	// sanity check.
	if frameWidth == 0 || frameHeight == 0 || frameWidth > 640 || frameHeight > 480 {
		return nil, false
	}
	dbg.Printf("frame dimensions: %dx%d", frameWidth, frameHeight)
	bounds := image.Rect(0, 0, frameWidth, frameHeight)
	dst := image.NewRGBA(bounds)
	data := frameContents[4:]
	drawPixel, total := pixelDrawer(dst, frameWidth, frameHeight)
	for pos := 0; pos < len(data); pos += 2 {
		cmd := binary.LittleEndian.Uint16(data[pos : pos+2])
		dbg.Printf("cmd: 0x%04X", cmd)
		if cmd == 0 {
			break
		}
		switch {
		case cmd&0x4000 != 0:
			// transparent lines.
			ySkip := int(cmd & 0xFFF)
			skip := ySkip * frameWidth
			for j := 0; j < skip; j++ {
				drawPixel(color.Transparent)
			}
		case cmd&0x1000 != 0:
			// regular pixels.
			npixels := int(cmd & 0xFFF)
			for j := 0; j < npixels; j++ {
				palIndex := data[pos]
				pos++
				drawPixel(pal[palIndex])
			}
		default:
			// transparent pixels.
			xSkip := int(cmd & 0xFFF)
			skip := xSkip
			for j := 0; j < skip; j++ {
				drawPixel(color.Transparent)
			}
		}
	}
	if *total > frameWidth*frameHeight {
		panic(fmt.Errorf("mismatch between total pixels drawn (%d) and expected image size (%dx%d = %d)", *total, frameWidth, frameHeight, frameWidth*frameHeight))
	}
	return dst, true
}

// pixelDrawer returns a function which may be invoked to incrementally set
// pixels; starting in the lower left corner, going from left to right, and then
// row by row from the bottom to the top of the image.
func pixelDrawer(dst draw.Image, w, h int) (func(color.Color), *int) {
	total := 0
	x, y := 0, 0
	return func(c color.Color) {
		// TODO: Remove sanity check once the zel decoder library has mature.
		if x < 0 || x >= w {
			panic(fmt.Sprintf("zel.pixelDrawer.drawPixel: invalid x; expected 0 <= x < %d, got x=%d", w, x))
		}
		if y < 0 || y >= h {
			panic(fmt.Sprintf("zel.pixelDrawer.drawPixel: invalid y; expected 0 <= y < %d, got y=%d", h, y))
		}
		total++
		dst.Set(x, y, c)
		x++
		if x >= w {
			x = 0
			y++
		}
	}, &total
}
