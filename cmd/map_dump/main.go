package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/term"
	"github.com/mewspring/pak/level/maps"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "map_dump:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("map_dump:")+" ", 0)
	// warn is a logger with the "map_dump:" prefix which logs warning messages
	// to standard error.
	warn = log.New(os.Stderr, term.RedBold("map_dump:")+" ", log.Lshortfile)
)

func usage() {
	const usage = "Usage: map_dump [OPTIONS]... FILE.map..."
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	// parse command line arguments.
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	// dump MAP files.
	for _, mapPath := range flag.Args() {
		if err := dumpMap(mapPath); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

// dumpMap dumps the given MAP file.
func dumpMap(mapPath string) error {
	m, err := maps.ParseFile(mapPath)
	if err != nil {
		return errors.WithStack(err)
	}
	pretty.Println("map:", m)
	return nil
}
