// ZEL file format
//
//    frame_offsets [...]uint32
//    frame_data    []byte

package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/mewkiz/pkg/imgutil"
	"github.com/mewkiz/pkg/pathutil"
	"github.com/mewkiz/pkg/term"
	"github.com/mewspring/pak/image/zel"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "zel_dump:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("zel_dump:")+" ", 0)
	// warn is a logger with the "zel_dump:" prefix which logs warning messages
	// to standard error.
	warn = log.New(os.Stderr, term.RedBold("zel_dump:")+" ", log.Lshortfile)
)

func usage() {
	const usage = "Usage: zel_dump [OPTIONS]... FILE.zel..."
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	// parse command line arguments.
	var palPath string
	flag.StringVar(&palPath, "pal", "", "palette path (256 RGBA colours)")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	// parse palette.
	pal, err := zel.ParsePal(palPath)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// dump ZEL image frames.
	for _, zelPath := range flag.Args() {
		if err := dumpZelImage(zelPath, pal); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

// dumpZelImage dumps the given ZEL image to the specified output directory.
func dumpZelImage(zelPath string, pal color.Palette) error {
	imgs, err := zel.DecodeAll(zelPath, pal)
	if err != nil {
		return errors.WithStack(err)
	}
	// create output directory.
	dstDir := pathutil.TrimExt(zelPath)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return errors.WithStack(err)
	}
	// output frames.
	for i, img := range imgs {
		pngName := fmt.Sprintf("frame_%04d.png", i)
		pngPath := filepath.Join(dstDir, pngName)
		dbg.Printf("creating %q", pngPath)
		if err := imgutil.WriteFile(pngPath, img); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
