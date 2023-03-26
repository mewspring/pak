// Package maps provides access to MAP files.
package maps

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"

	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
)

var (
	// dbg is a logger with the "map:" prefix which logs debug messages to
	// standard error.
	dbg = log.New(os.Stderr, term.MagentaBold("map:")+" ", 0)
	// warn is a logger with the "map:" prefix which logs warning messages
	// to standard error.
	warn = log.New(os.Stderr, term.RedBold("map:")+" ", log.Lshortfile)
)

// signature specifies the file format signature of MAP files.
const signature = "MAP\x00"

// Map holds the contents of a MAP file.
type Map struct {
	// File format signature "MAP\x00"
	Magic           [4]byte
	Unused0004      uint32
	RenderWithLight uint8
	// Walls tileset ID.
	//
	// Subfile index of X/base_walls_tileset/base_walls_NNN.zel
	BaseWallsTilesetID uint32
	// Collisions of the map.
	SolidMap [128][128]uint8
	// Floor tile frame indices of the map.
	//
	//    frame >= 0: use base floors tileset (X/base_floors_tileset.zel)
	//    else:       use tileset type 2      (X/tilesets/tileset_NNN_floors.zel)
	FloorFrameMap [128][128]uint16
	// Tileset type 0 holds stairs and mountains
	//
	//	X/tilesets/tileset_NNN_mountains_and_stairs.zel
	//
	// ntileset0Elems uint32
	Tileset0Elems []MapTile // len: ntileset0Elems
	// Tileset type 4 holds shadows
	//
	//	X/tilesets/tileset_NNN_shadows.zel
	//
	// ntileset4Elems uint32
	Tileset4Elems []MapTile // len: ntileset4Elems
	// Tileset type 1 holds walls and buildings
	//
	//	X/tilesets/tileset_NNN_walls_and_buildings.zel
	//
	// ntileset1Elems uint32 // in range [0, 4096)
	Tileset1Elems []MapTile2 // len: ntileset1Elems
	// Tileset type 3 holds objects
	//
	//	X/tilesets/tileset_NNN_objects.zel
	//
	// ntileset3Elems uint32 // in range [0, 4096)
	Tileset3Elems []MapTile2 // len: ntileset3Elems
	// Base walls tileset (of X subarchive 4)
	//
	//	X/base_walls_tileset/base_walls_NNN.zel
	//
	// ntilesetWallsElems uint32 // in range [0, 4096)
	TilesetWallsElems []MapTile2 // len: ntilesetWallsElems
}

// MapTile specifies the tileset frame index and map coordinate of a map tile.
type MapTile struct {
	Frame uint16
	_     [2]byte // padding
	X     uint32
	Y     uint32
	_     [8]byte // padding
}

// MapTile2 specifies the tileset frame index and map coordinate of a map tile.
type MapTile2 struct {
	Frame uint16
	X     uint8
	Y     uint8
}

// ParseFile parses the given MAP file.
func ParseFile(mapPath string) (*Map, error) {
	dbg.Printf("parsing %q", mapPath)
	buf, err := ioutil.ReadFile(mapPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r := bytes.NewReader(buf)
	m := &Map{}
	if err := binary.Read(r, binary.LittleEndian, &m.Magic); err != nil {
		return nil, errors.WithStack(err)
	}
	magic := string(m.Magic[:])
	if magic != signature {
		return nil, errors.Errorf("invalid MAP signature of %q; expected %q, got %q", mapPath, signature, magic)
	}
	//dbg.Println("magic:", magic)
	if err := binary.Read(r, binary.LittleEndian, &m.Unused0004); err != nil {
		return nil, errors.WithStack(err)
	}
	//dbg.Printf("m.Unused0004: 0x%08X", m.Unused0004)
	if err := binary.Read(r, binary.LittleEndian, &m.RenderWithLight); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.RenderWithLight: 0x%02X", m.RenderWithLight)
	if err := binary.Read(r, binary.LittleEndian, &m.BaseWallsTilesetID); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.BaseWallsTilesetID: 0x%02X", m.BaseWallsTilesetID)
	if err := binary.Read(r, binary.LittleEndian, &m.SolidMap); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.SolidMap:\n%v", m.SolidMap)
	// Floors.
	if err := binary.Read(r, binary.LittleEndian, &m.FloorFrameMap); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.FloorFrameMap:\n%v", m.FloorFrameMap)
	// Tileset 0 (stairs and mountains).
	var ntileset0Elems uint32
	if err := binary.Read(r, binary.LittleEndian, &ntileset0Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("ntileset0Elems:", ntileset0Elems)
	m.Tileset0Elems = make([]MapTile, int(ntileset0Elems))
	if err := binary.Read(r, binary.LittleEndian, &m.Tileset0Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.Tileset0Elems (stairs and mountains):\n%v", m.Tileset0Elems)
	// Tileset 4 (shadows).
	var ntileset4Elems uint32
	if err := binary.Read(r, binary.LittleEndian, &ntileset4Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("ntileset4Elems:", ntileset4Elems)
	m.Tileset4Elems = make([]MapTile, int(ntileset4Elems))
	if err := binary.Read(r, binary.LittleEndian, &m.Tileset4Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.Tileset4Elems (shadows):\n%v", m.Tileset4Elems)
	// Tileset 1 (walls and buildings).
	var ntileset1Elems uint32
	if err := binary.Read(r, binary.LittleEndian, &ntileset1Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("ntileset1Elems:", ntileset1Elems)
	m.Tileset1Elems = make([]MapTile2, int(ntileset1Elems))
	if err := binary.Read(r, binary.LittleEndian, &m.Tileset1Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.Tileset1Elems (walls and buildings):\n%v", m.Tileset1Elems)
	// Tileset 3 (objects).
	var ntileset3Elems uint32
	if err := binary.Read(r, binary.LittleEndian, &ntileset3Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("ntileset3Elems:", ntileset3Elems)
	m.Tileset3Elems = make([]MapTile2, int(ntileset3Elems))
	if err := binary.Read(r, binary.LittleEndian, &m.Tileset3Elems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.Tileset3Elems (objects):\n%v", m.Tileset3Elems)
	// Base walls.
	var ntilesetWallsElems uint32
	if err := binary.Read(r, binary.LittleEndian, &ntilesetWallsElems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("ntilesetWallsElems:", ntilesetWallsElems)
	m.TilesetWallsElems = make([]MapTile2, int(ntilesetWallsElems))
	if err := binary.Read(r, binary.LittleEndian, &m.TilesetWallsElems); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("m.TilesetWallsElems:\n%v", m.TilesetWallsElems)
	return m, nil
}
