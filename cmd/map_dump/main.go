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
	baseFloorBaseTileID     = 0 + 1       // X/base_floor_tileset.zel
	tileset1FloorBaseTileID = 100_000 + 1 // X/tilesets/tileset_NNN_floors.zel
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
		dbg.Println("positive floorFrame:", floorFrame)
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
		dbg.Println("negative floorFrame:", floorFrame)
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

// addTilesets add all tilesets to the given TMX map.
func addTilesets(tmxMap *tmx.Map) {
	// Add base floors tileset.
	baseFloorsTilesetInfo := &TilesetInfo{
		TilesetName:       "base_floors",
		TilesetWidth:      768,  // TODO: infer from image.
		TilesetHeight:     2272, // TODO: infer from image.
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        baseFloorBaseTileID,
	}
	addTileset(tmxMap, baseFloorsTilesetInfo)
	// Add tileset NNN floors tileset.
	floorsTilesetInfo := &TilesetInfo{
		TilesetName:       "tileset_1/floors", // TODO: add tileset_NNN/ support.
		TilesetWidth:      768,                // TODO: infer from image.
		TilesetHeight:     800,                // TODO: infer from image.
		TilesetTileWidth:  mapTileWidth,
		TilesetTileHeight: mapTileHeight,
		BaseTileID:        tileset1FloorBaseTileID,
	}
	addTileset(tmxMap, floorsTilesetInfo)
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
