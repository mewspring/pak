package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Noofbiz/tmx"
	"github.com/mewkiz/pkg/pathutil"
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

const (
	// Map width in number of tiles.
	mapWidth = 128
	// Map height in number of tiles.
	mapHeight = 128
	// Tile width (in pixels) of each tile on the map.
	mapTileWidth = 64
	// Tile height (in pixels) of each tile on the map.
	mapTileHeight = 32
)

const (
	baseFloorBaseTileID      = 0 + 1          // X/base_floors_tileset.zel
	tileset1FloorBaseTileID  = 1*100_000 + 1  // X/tilesets/tileset_1_floors.zel
	tileset2FloorBaseTileID  = 2*100_000 + 1  // X/tilesets/tileset_2_floors.zel
	tileset3FloorBaseTileID  = 3*100_000 + 1  // X/tilesets/tileset_3_floors.zel
	tileset4FloorBaseTileID  = 4*100_000 + 1  // X/tilesets/tileset_4_floors.zel
	tileset5FloorBaseTileID  = 5*100_000 + 1  // X/tilesets/tileset_5_floors.zel
	tileset6FloorBaseTileID  = 6*100_000 + 1  // X/tilesets/tileset_6_floors.zel
	tileset7FloorBaseTileID  = 7*100_000 + 1  // X/tilesets/tileset_7_floors.zel
	tileset8FloorBaseTileID  = 8*100_000 + 1  // X/tilesets/tileset_8_floors.zel
	tileset9FloorBaseTileID  = 9*100_000 + 1  // X/tilesets/tileset_9_floors.zel
	tileset10FloorBaseTileID = 10*100_000 + 1 // X/tilesets/tileset_10_floors.zel
	tileset11FloorBaseTileID = 11*100_000 + 1 // X/tilesets/tileset_11_floors.zel
	tileset12FloorBaseTileID = 12*100_000 + 1 // X/tilesets/tileset_12_floors.zel
	tileset13FloorBaseTileID = 13*100_000 + 1 // X/tilesets/tileset_13_floors.zel
	tileset14FloorBaseTileID = 14*100_000 + 1 // X/tilesets/tileset_14_floors.zel
	tileset15FloorBaseTileID = 15*100_000 + 1 // X/tilesets/tileset_15_floors.zel
	tileset16FloorBaseTileID = 16*100_000 + 1 // X/tilesets/tileset_16_floors.zel
	tileset17FloorBaseTileID = 17*100_000 + 1 // X/tilesets/tileset_17_floors.zel
)

const outputDir = "_assets_"

// dumpMap dumps the given MAP file.
func dumpMap(mapPath string) error {
	// Parse MAP file.
	m, err := maps.ParseFile(mapPath)
	if err != nil {
		return errors.WithStack(err)
	}
	// Convert MAP file to TMX format.
	tmxMap := convertMapToTmx(m)
	// Output TMX map.
	// Store TMX file to output directory.
	dstDir := filepath.Join(outputDir, "maps")
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	mapName := pathutil.FileName(mapPath)
	tmxName := fmt.Sprintf("%s.tmx", mapName)
	tmxPath := filepath.Join(dstDir, tmxName)
	if err := dumpTmx(tmxPath, tmxMap); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// convertMapToTmx converts the given MAP file to TMX format.
func convertMapToTmx(m *maps.Map) *tmx.Map {
	//pretty.Println("map:", m)
	//fmt.Println("solid map:")
	//for y := 0; y < 128; y++ {
	//	for x := 0; x < 128; x++ {
	//		fmt.Printf("%d", m.SolidMap[y][x])
	//	}
	//	fmt.Println()
	//}
	// Create TMX map.
	tmxMap := &tmx.Map{
		Orientation: "isometric",
		Width:       mapWidth,
		Height:      mapHeight,
		TileWidth:   mapTileWidth,
		TileHeight:  mapTileHeight,
	}
	// Add TMX tilesets.
	addTilesets(tmxMap)
	// Add layers.
	addLayers(tmxMap, m)
	return tmxMap
}

// addLayers converts MAP layers to TMX format.
func addLayers(tmxMap *tmx.Map, m *maps.Map) {
	// Base floor layer.
	const baseFloorLayerName = "base_floor"
	baseFloorTiledTileIDAt := func(x, y int) int {
		floorFrame := int(int16(m.FloorFrameMap[y][x]))
		// base floor tile IDs are positive.
		if floorFrame < 0 {
			return 0 // early return; ignore negative (tileset type 2, floors).
		}
		//dbg.Println("positive floorFrame:", floorFrame)
		const base = baseFloorBaseTileID
		tiledTileID := base + floorFrame
		return tiledTileID
	}
	addLayer(tmxMap, baseFloorLayerName, baseFloorTiledTileIDAt)

	// Floor layer.
	const floorLayerName = "floor"
	floorTiledTileIDAt := func(x, y int) int {
		floorFrame := int(int16(m.FloorFrameMap[y][x]))
		// tileset type 2 (floors) tile IDs are negative.
		if floorFrame >= 0 {
			return 0 // early return; ignore positive (base floors) and 0.
		}
		floorFrame &= 0x7FFF
		//dbg.Println("negative floorFrame:", floorFrame)
		const base = tileset1FloorBaseTileID // TODO: add support for tileset_NNN/ type.
		tiledTileID := base + floorFrame
		return tiledTileID
	}
	addLayer(tmxMap, floorLayerName, floorTiledTileIDAt)
}

// addLayer adds the layer as specified to the given TMX map.
func addLayer(tmxMap *tmx.Map, tilesLayerName string, tiledTileIDAt func(x, y int) int) {
	inner := &strings.Builder{}
	inner.WriteString("\n")
	for y := 0; y < mapWidth; y++ {
		for x := 0; x < mapHeight; x++ {
			tiledTileID := tiledTileIDAt(x, y)
			inner.WriteString(strconv.Itoa(tiledTileID))
			if !(x == mapWidth-1 && y == mapHeight-1) {
				inner.WriteString(",")
			}
		}
		inner.WriteString("\n")
	}
	data := tmx.Data{
		Encoding: "csv",
		Inner:    inner.String(),
	}
	// Add tiles.
	tilesLayer := tmx.Layer{
		Name:    tilesLayerName,
		Data:    []tmx.Data{data},
		Width:   mapWidth,
		Height:  mapHeight,
		Opacity: 1.0,
		Visible: 1,
	}
	tmxMap.Layers = append(tmxMap.Layers, tilesLayer)
}

// TilesetInfo specifies the name, dimensions and base tile ID of a given
// tileset.
type TilesetInfo struct {
	TilesetName       string
	TilesetWidth      int
	TilesetHeight     int
	TilesetTileWidth  int
	TilesetTileHeight int
	BaseTileID        int
}

// floorsTilesetInfos specifies the tileset info of
// `X/tilesets/tileset_NNN_floors.zel` floor tilesets
var floorsTilesetInfos = []*TilesetInfo{
	{
		TilesetName:       "tileset_1/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset1FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_2/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset2FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_3/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset3FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_4/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset4FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_5/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset5FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_6/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset6FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_7/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset7FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_8/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset8FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_9/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset9FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_10/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset10FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_11/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset11FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_12/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset12FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_13/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset13FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_14/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset14FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_15/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset15FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_16/floors",
		TilesetWidth:      768,
		TilesetHeight:     320,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset16FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_17/floors",
		TilesetWidth:      768,
		TilesetHeight:     64,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset17FloorBaseTileID,
	},
}

// addTilesets add all tilesets to the given TMX map.
func addTilesets(tmxMap *tmx.Map) {
	// Add base floors tileset.
	baseFloorsTilesetInfo := &TilesetInfo{
		TilesetName:       "base_floors",
		TilesetWidth:      768,
		TilesetHeight:     2272,
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        baseFloorBaseTileID,
	}
	addTileset(tmxMap, baseFloorsTilesetInfo)
	// Add tileset NNN floors tilesets.
	for _, floorsTilesetInfo := range floorsTilesetInfos {
		addTileset(tmxMap, floorsTilesetInfo)
	}
}

// addTileset add the tileset as specified to the given TMX map.
func addTileset(tmxMap *tmx.Map, tilesetInfo *TilesetInfo) {
	tilesetPngName := fmt.Sprintf("%s.png", tilesetInfo.TilesetName)
	tilesetImg := tmx.Image{
		// Traverse up to the root of the "_assets_" directory.
		Source: filepath.Join("..", "tilesets", tilesetPngName),
		Width:  float64(tilesetInfo.TilesetWidth),
		Height: float64(tilesetInfo.TilesetHeight),
	}
	tileset := tmx.Tileset{
		FirstGID:   uint32(tilesetInfo.BaseTileID),
		Name:       tilesetInfo.TilesetName,
		TileWidth:  tilesetInfo.TilesetTileWidth,
		TileHeight: tilesetInfo.TilesetTileHeight,
		Image:      []tmx.Image{tilesetImg},
	}
	tmxMap.Tilesets = append(tmxMap.Tilesets, tileset)
}

// dumpTmx stores the given TMX map to the specified TMX path.
func dumpTmx(tmxPath string, tmxMap *tmx.Map) error {
	dbg.Printf("Creating %q.", tmxPath)
	buf, err := xml.MarshalIndent(tmxMap, "", "\t")
	if err != nil {
		return errors.WithStack(err)
	}

	// TODO: figure out why XML marshal uses title cased <Map> tag. Tiled
	// requires lower case.
	buf = bytes.ReplaceAll(buf, []byte("<Map"), []byte("<map"))
	buf = bytes.ReplaceAll(buf, []byte("</Map"), []byte("</map"))

	if err := ioutil.WriteFile(tmxPath, buf, 0644); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
