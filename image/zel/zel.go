// Package zel provides access to ZEL image files.
package zel

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/imgutil"
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
func DecodeAll(zelPath string, pal color.Palette) (imgs []image.Image, err error) {
	var (
		curFrame int
		nframes  int
	)
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("recovered panic in zel.DecodeAll for frame (%d/%d) of %q: %+v", curFrame, nframes, zelPath, e)
		}
	}()
	buf, err := ioutil.ReadFile(zelPath)
	if err != nil {
		return imgs, errors.WithStack(err)
	}
	dbg.Printf("parsing %q", zelPath)
	// parse ZEL header.
	zelHdrSize := int(binary.LittleEndian.Uint32(buf[0:4]))
	r := bytes.NewReader(buf)
	zelHdrReader := io.NewSectionReader(r, 0, int64(zelHdrSize))
	frameOffsetsLen := zelHdrSize / 4
	frameOffsets := make([]uint32, frameOffsetsLen)
	if err := binary.Read(zelHdrReader, binary.LittleEndian, &frameOffsets); err != nil {
		return nil, errors.WithStack(err)
	}
	nframes = len(frameOffsets) - 1
	if len(buf) != int(frameOffsets[nframes]) {
		pretty.Println("frameOffsets:", frameOffsets)
		panic(fmt.Errorf("mismatch between frameOffsets[%d]=%d and len(buf)=%d", nframes, frameOffsets[nframes], len(buf)))
	}
	// output ZEL frames.
	for curFrame = 0; curFrame < nframes; curFrame++ {
		frameStartOffset := frameOffsets[curFrame]
		frameEndOffset := frameOffsets[curFrame+1]
		frameContents := buf[frameStartOffset:frameEndOffset:frameEndOffset]
		type4 := isType4(zelPath)
		img, ok := parseFrame(frameContents, pal, type4)
		if !ok {
			//warn.Printf("skipping invalid frame (%d/%d) of %q", curFrame, nframes, zelPath)
			//continue // skip
			return imgs, errors.Errorf("unable to decode frame (%d/%d) of %q", curFrame, nframes, zelPath)
		}
		imgs = append(imgs, img)
	}
	return imgs, errors.WithStack(err)
}

// parseFrame parses the given ZEL frame contents.
func parseFrame(frameContents []byte, pal color.Palette, type4 bool) (image.Image, bool) {
	// parse ZEL frame.
	if len(frameContents) == 0 {
		warn.Printf("empty frame")
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), true // dummy 1x1 image used for empty frames
	}
	frameWidth := int(binary.LittleEndian.Uint16(frameContents[0:2]))
	frameHeight := int(binary.LittleEndian.Uint16(frameContents[2:4]))
	// sanity check.
	// NOTE: 650 is a valid width of `archive_0012/archive_0005/frame_0000.png`.
	// NOTE: 640 is a valid height of `archive_0012/archive_0001/frame_0165.png`.
	// NOTE: 1037 is a valid width of `archive_0012/archive_0040/frame_0034.png`.
	if frameWidth == 0 || frameHeight == 0 || frameWidth > 1280 || frameHeight > 1280 {
		warn.Printf("sanity check failed; frameWidth=%d, frameHeight=%d\n%s", frameWidth, frameHeight, hex.Dump(frameContents))
		return nil, false
	}
	dbg.Printf("frame dimensions: %dx%d", frameWidth, frameHeight)
	bounds := image.Rect(0, 0, frameWidth, frameHeight)
	dst := image.NewRGBA(bounds)

	// TODO: remove partial_img.png output after broken ZEL images have been
	// patched. Used now for debugging.
	defer func() {
		if e := recover(); e != nil {
			if err := imgutil.WriteFile("_dump_/partial_img.png", dst); err != nil {
				panic(err)
			}
			warn.Printf("recovered from %+v", e)
		}
	}()

	data := frameContents[4:]
	drawPixel, total := pixelDrawer(dst, frameWidth, frameHeight)
	for pos := 0; pos < len(data); {
		cmd := binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2
		//dbg.Printf("cmd: 0x%04X", cmd)
		if cmd == 0 {
			if pos < len(data) {
				warn.Printf("unprocessed frame contents:\n%s", hex.Dump(data[pos:]))
			}
			break
		}
		switch {
		case cmd&0x4000 != 0:
			// transparent lines.
			ySkip := int(cmd & 0xFFF)
			//dbg.Printf("   transparent lines (ySkip=%d)", ySkip)
			if ySkip > frameHeight {
				panic(fmt.Errorf("invalid ySkip (%d); exceeds frame width (%d)", ySkip, frameWidth))
			}
			skip := ySkip * frameWidth
			for j := 0; j < skip; j++ {
				drawPixel(color.Transparent)
			}
		case cmd&0x1000 != 0:
			// regular pixels.
			npixels := int(cmd & 0xFFF)
			switch {
			case type4:
				// Tileset shadows (using constant palette index 8).
				//
				//    "X/tilesets/archive_NNNN.zel" where (NNNN%4 == 0)
				//dbg.Printf("   constant pixels (npixels=%d)", npixels)
				if npixels > frameWidth {
					panic(fmt.Errorf("invalid npixels (%d); exceeds frame width (%d)", npixels, frameWidth))
				}
				for j := 0; j < npixels; j++ {
					const palIndex = 8
					//dbg.Printf("      constant pixel 0x%02X", palIndex)
					drawPixel(pal[palIndex])
				}
			default:
				//dbg.Printf("   regular pixels (npixels=%d)", npixels)
				if npixels > frameWidth {
					panic(fmt.Errorf("invalid npixels (%d); exceeds frame width (%d)", npixels, frameWidth))
				}
				for j := 0; j < npixels; j++ {
					palIndex := data[pos]
					//dbg.Printf("      regular pixel 0x%02X", palIndex)
					pos++
					drawPixel(pal[palIndex])
				}
			}
		default:
			// transparent pixels.
			xSkip := int(cmd & 0xFFF)
			//dbg.Printf("   transparent pixels (xSkip=%d)", xSkip)
			if xSkip > frameWidth {
				panic(fmt.Errorf("invalid xSkip (%d); exceeds frame width (%d)", xSkip, frameWidth))
			}
			skip := xSkip
			for j := 0; j < skip; j++ {
				drawPixel(color.Transparent)
			}
		}
		if cmd&0x8000 != 0 {
			// clear line.
			skip := (frameWidth - *total) % frameWidth // TODO: double-check that this is correct
			//dbg.Printf("   clear line (skip=%d)", skip)
			if skip != 0 {
				panic(fmt.Errorf("unexpected clear line skip; expected 0, got %d", skip))
			}
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

// isType4 reports whether the given ZEL is a type 4 tileset ZEL image (used for
// tileset shadows).
func isType4(zelPath string) bool {
	zelPath = strings.ReplaceAll(zelPath, `\`, "/")
	const rootDir = "X/"
	pos := strings.Index(zelPath, rootDir)
	if pos == -1 {
		panic(fmt.Errorf("unable to find root directory %q in %q", rootDir, zelPath))
	}
	zelPath = zelPath[pos:]
	return isType4TilesetZel[zelPath]
}

// isType4TilesetZel reports whether the given ZEL is a type 4 tileset ZEL
// image (used for tileset shadows).
var isType4TilesetZel = map[string]bool{
	"X/tilesets/tileset_1_shadows.zel":  true,
	"X/tilesets/tileset_2_shadows.zel":  true,
	"X/tilesets/tileset_3_shadows.zel":  true,
	"X/tilesets/tileset_4_shadows.zel":  true,
	"X/tilesets/tileset_5_shadows.zel":  true,
	"X/tilesets/tileset_6_shadows.zel":  true,
	"X/tilesets/tileset_7_shadows.zel":  true,
	"X/tilesets/tileset_8_shadows.zel":  true,
	"X/tilesets/tileset_9_shadows.zel":  true,
	"X/tilesets/tileset_10_shadows.zel": true,
	"X/tilesets/tileset_11_shadows.zel": true,
	"X/tilesets/tileset_12_shadows.zel": true,
	"X/tilesets/tileset_13_shadows.zel": true,
	"X/tilesets/tileset_14_shadows.zel": true,
	"X/tilesets/tileset_15_shadows.zel": true,
	"X/tilesets/tileset_16_shadows.zel": true,
	"X/tilesets/tileset_17_shadows.zel": true,
}
