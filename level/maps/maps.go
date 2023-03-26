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
	// nmountains uint32
	Mountains []MapTile // len: nmountains
	// Tileset type 4 holds shadows
	//
	//	X/tilesets/tileset_NNN_shadows.zel
	//
	// nshadows uint32
	Shadows []MapTile // len: nshadows
	// Tileset type 1 holds walls and buildings
	//
	//	X/tilesets/tileset_NNN_walls_and_buildings.zel
	//
	// nbuildings uint32 // in range [0, 4096)
	Buildings []MapTile2 // len: nbuildings
	// Tileset type 3 holds objects
	//
	//	X/tilesets/tileset_NNN_objects.zel
	//
	// nobjects uint32 // in range [0, 4096)
	Objects []MapTile2 // len: nobjects
	// Base walls tileset (of X subarchive 4)
	//
	//	X/base_walls_tileset/base_walls_NNN.zel
	//
	// ntilesetWallsElems uint32 // in range [0, 4096)
	BaseWalls []MapTile2 // len: ntilesetWallsElems
}

// MapTile specifies the tileset frame index and map coordinate of a map tile.
type MapTile struct {
	Frame uint16
	_     [2]byte // padding
	X     int32
	Y     int32
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
	//dbg.Printf("m.SolidMap:\n%v", m.SolidMap)
	// Floors.
	if err := binary.Read(r, binary.LittleEndian, &m.FloorFrameMap); err != nil {
		return nil, errors.WithStack(err)
	}
	//dbg.Printf("m.FloorFrameMap:\n%v", m.FloorFrameMap)
	// Tileset 0 (stairs and mountains).
	var nmountains uint32
	if err := binary.Read(r, binary.LittleEndian, &nmountains); err != nil {
		return nil, errors.WithStack(err)
	}
	m.Mountains = make([]MapTile, int(nmountains))
	if err := binary.Read(r, binary.LittleEndian, &m.Mountains); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.Mountains (stairs and mountains):")
	for _, mountain := range m.Mountains {
		dbg.Println("   mountain:", mountain)
	}
	// Tileset 4 (shadows).
	var nshadows uint32
	if err := binary.Read(r, binary.LittleEndian, &nshadows); err != nil {
		return nil, errors.WithStack(err)
	}
	m.Shadows = make([]MapTile, int(nshadows))
	if err := binary.Read(r, binary.LittleEndian, &m.Shadows); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.Shadows (shadows):")
	for _, shadow := range m.Shadows {
		dbg.Println("   shadow:", shadow)
	}
	// Tileset 1 (walls and buildings).
	var nbuildings uint32
	if err := binary.Read(r, binary.LittleEndian, &nbuildings); err != nil {
		return nil, errors.WithStack(err)
	}
	m.Buildings = make([]MapTile2, int(nbuildings))
	if err := binary.Read(r, binary.LittleEndian, &m.Buildings); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.Buildings (walls and buildings):")
	for _, building := range m.Buildings {
		dbg.Println("   building:", building)
	}
	// Tileset 3 (objects).
	var nobjects uint32
	if err := binary.Read(r, binary.LittleEndian, &nobjects); err != nil {
		return nil, errors.WithStack(err)
	}
	m.Objects = make([]MapTile2, int(nobjects))
	if err := binary.Read(r, binary.LittleEndian, &m.Objects); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.Objects (objects):")
	for _, object := range m.Objects {
		dbg.Println("   object:", object)
	}
	// Base walls.
	var ntilesetWallsElems uint32
	if err := binary.Read(r, binary.LittleEndian, &ntilesetWallsElems); err != nil {
		return nil, errors.WithStack(err)
	}
	m.BaseWalls = make([]MapTile2, int(ntilesetWallsElems))
	if err := binary.Read(r, binary.LittleEndian, &m.BaseWalls); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.BaseWalls:")
	for _, baseWall := range m.BaseWalls {
		dbg.Println("   baseWall:", baseWall)
	}
	return m, nil
}
