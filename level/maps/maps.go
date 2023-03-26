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
	Mountains []MapOverlay // len: nmountains
	// Tileset type 4 holds shadows.
	//
	//	X/tilesets/tileset_NNN_shadows.zel
	//
	// nshadows uint32
	Shadows []MapOverlay // len: nshadows
	// Tileset type 1 holds buildings.
	//
	//	X/tilesets/tileset_NNN_buildings.zel
	//
	// nbuildings uint32 // in range [0, 4096)
	Buildings []MapTile // len: nbuildings
	// Tileset type 3 holds objects.
	//
	//	X/tilesets/tileset_NNN_objects.zel
	//
	// nobjects uint32 // in range [0, 4096)
	Objects []MapTile // len: nobjects
	// Base walls tilesets (of X subarchive 4).
	//
	//	X/base_walls_tileset/base_walls_NNN.zel
	//
	// nbaseWalls uint32 // in range [0, 4096)
	BaseWalls []MapTile // len: nbaseWalls
}

// MapOverlay specifies the tileset frame index and screen offset of a map
// overlay.
type MapOverlay struct {
	// Tileset frame index.
	Frame uint16
	_     [2]byte // padding
	// (X,Y)-screen offset in pixels.
	X int32
	Y int32
	_ [8]byte // padding
}

// MapTile specifies the tileset frame index and map coordinate of a map tile.
type MapTile struct {
	// Tileset frame index.
	Frame uint16
	// (X,Y) map coordinate.
	X uint8
	Y uint8
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
	m.Mountains = make([]MapOverlay, int(nmountains))
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
	m.Shadows = make([]MapOverlay, int(nshadows))
	if err := binary.Read(r, binary.LittleEndian, &m.Shadows); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.Shadows:")
	for _, shadow := range m.Shadows {
		dbg.Println("   shadow:", shadow)
	}
	// Tileset 1 (buildings).
	var nbuildings uint32
	if err := binary.Read(r, binary.LittleEndian, &nbuildings); err != nil {
		return nil, errors.WithStack(err)
	}
	m.Buildings = make([]MapTile, int(nbuildings))
	if err := binary.Read(r, binary.LittleEndian, &m.Buildings); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.Buildings:")
	for _, building := range m.Buildings {
		dbg.Println("   building:", building)
	}
	// Tileset 3 (objects).
	var nobjects uint32
	if err := binary.Read(r, binary.LittleEndian, &nobjects); err != nil {
		return nil, errors.WithStack(err)
	}
	m.Objects = make([]MapTile, int(nobjects))
	if err := binary.Read(r, binary.LittleEndian, &m.Objects); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.Objects:")
	for _, object := range m.Objects {
		dbg.Println("   object:", object)
	}
	// Base walls.
	var nbaseWalls uint32
	if err := binary.Read(r, binary.LittleEndian, &nbaseWalls); err != nil {
		return nil, errors.WithStack(err)
	}
	m.BaseWalls = make([]MapTile, int(nbaseWalls))
	if err := binary.Read(r, binary.LittleEndian, &m.BaseWalls); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Println("m.BaseWalls:")
	for _, baseWall := range m.BaseWalls {
		dbg.Println("   baseWall:", baseWall)
	}
	return m, nil
}
