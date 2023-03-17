package zel

import (
	"fmt"
	"image/color"
	"image/color/palette"
	"io/ioutil"

	"github.com/pkg/errors"
)

// ParsePal parses the given RGBA palette.
func ParsePal(palPath string) (color.Palette, error) {
	if len(palPath) == 0 {
		// use hardcoded fallback palette.
		warn.Printf("using fallback palette; use -pal to specify palette")
		return palette.Plan9, nil
	}
	buf, err := ioutil.ReadFile(palPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	const ncolors = 256
	if len(buf) != ncolors*4 {
		panic(fmt.Errorf("invalid palette length; expected 256*4, got %d", len(buf)))
	}
	pal := make(color.Palette, ncolors)
	for i := range pal {
		c := color.RGBA{
			R: buf[i*4+0],
			G: buf[i*4+1],
			B: buf[i*4+2],
			A: 0xFF, // buf[i*4+3]
		}
		pal[i] = c
	}
	return pal, nil
}
