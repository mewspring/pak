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
	// base floors.
	baseFloorBaseTileID = 0 + 1 // X/base_floors_tileset.zel
	// base walls.
	baseWalls1BaseTileID = 200_000*0 + 1*10_000 + 1 // X/base_walls/base_walls_tileset/base_walls_1.zel
	baseWalls2BaseTileID = 200_000*0 + 2*10_000 + 1 // X/base_walls/base_walls_tileset/base_walls_2.zel
	baseWalls3BaseTileID = 200_000*0 + 3*10_000 + 1 // X/base_walls/base_walls_tileset/base_walls_3.zel
	baseWalls4BaseTileID = 200_000*0 + 4*10_000 + 1 // X/base_walls/base_walls_tileset/base_walls_4.zel
	baseWalls5BaseTileID = 200_000*0 + 5*10_000 + 1 // X/base_walls/base_walls_tileset/base_walls_5.zel
	baseWalls6BaseTileID = 200_000*0 + 6*10_000 + 1 // X/base_walls/base_walls_tileset/base_walls_6.zel
	baseWalls7BaseTileID = 200_000*0 + 7*10_000 + 1 // X/base_walls/base_walls_tileset/base_walls_7.zel
	// floors.
	tileset1FloorBaseTileID  = 200_000*1 + 1*10_000 + 1  // X/tilesets/tileset_1_floors.zel
	tileset2FloorBaseTileID  = 200_000*1 + 2*10_000 + 1  // X/tilesets/tileset_2_floors.zel
	tileset3FloorBaseTileID  = 200_000*1 + 3*10_000 + 1  // X/tilesets/tileset_3_floors.zel
	tileset4FloorBaseTileID  = 200_000*1 + 4*10_000 + 1  // X/tilesets/tileset_4_floors.zel
	tileset5FloorBaseTileID  = 200_000*1 + 5*10_000 + 1  // X/tilesets/tileset_5_floors.zel
	tileset6FloorBaseTileID  = 200_000*1 + 6*10_000 + 1  // X/tilesets/tileset_6_floors.zel
	tileset7FloorBaseTileID  = 200_000*1 + 7*10_000 + 1  // X/tilesets/tileset_7_floors.zel
	tileset8FloorBaseTileID  = 200_000*1 + 8*10_000 + 1  // X/tilesets/tileset_8_floors.zel
	tileset9FloorBaseTileID  = 200_000*1 + 9*10_000 + 1  // X/tilesets/tileset_9_floors.zel
	tileset10FloorBaseTileID = 200_000*1 + 10*10_000 + 1 // X/tilesets/tileset_10_floors.zel
	tileset11FloorBaseTileID = 200_000*1 + 11*10_000 + 1 // X/tilesets/tileset_11_floors.zel
	tileset12FloorBaseTileID = 200_000*1 + 12*10_000 + 1 // X/tilesets/tileset_12_floors.zel
	tileset13FloorBaseTileID = 200_000*1 + 13*10_000 + 1 // X/tilesets/tileset_13_floors.zel
	tileset14FloorBaseTileID = 200_000*1 + 14*10_000 + 1 // X/tilesets/tileset_14_floors.zel
	tileset15FloorBaseTileID = 200_000*1 + 15*10_000 + 1 // X/tilesets/tileset_15_floors.zel
	tileset16FloorBaseTileID = 200_000*1 + 16*10_000 + 1 // X/tilesets/tileset_16_floors.zel
	tileset17FloorBaseTileID = 200_000*1 + 17*10_000 + 1 // X/tilesets/tileset_17_floors.zel
	// objects.
	tileset1ObjectBaseTileID  = 200_000*2 + 1*10_000 + 1  // X/tilesets/tileset_1_objects.zel
	tileset2ObjectBaseTileID  = 200_000*2 + 2*10_000 + 1  // X/tilesets/tileset_2_objects.zel
	tileset3ObjectBaseTileID  = 200_000*2 + 3*10_000 + 1  // X/tilesets/tileset_3_objects.zel
	tileset4ObjectBaseTileID  = 200_000*2 + 4*10_000 + 1  // X/tilesets/tileset_4_objects.zel
	tileset5ObjectBaseTileID  = 200_000*2 + 5*10_000 + 1  // X/tilesets/tileset_5_objects.zel
	tileset6ObjectBaseTileID  = 200_000*2 + 6*10_000 + 1  // X/tilesets/tileset_6_objects.zel
	tileset7ObjectBaseTileID  = 200_000*2 + 7*10_000 + 1  // X/tilesets/tileset_7_objects.zel
	tileset8ObjectBaseTileID  = 200_000*2 + 8*10_000 + 1  // X/tilesets/tileset_8_objects.zel
	tileset9ObjectBaseTileID  = 200_000*2 + 9*10_000 + 1  // X/tilesets/tileset_9_objects.zel
	tileset10ObjectBaseTileID = 200_000*2 + 10*10_000 + 1 // X/tilesets/tileset_10_objects.zel
	tileset11ObjectBaseTileID = 200_000*2 + 11*10_000 + 1 // X/tilesets/tileset_11_objects.zel
	tileset12ObjectBaseTileID = 200_000*2 + 12*10_000 + 1 // X/tilesets/tileset_12_objects.zel
	tileset13ObjectBaseTileID = 200_000*2 + 13*10_000 + 1 // X/tilesets/tileset_13_objects.zel
	tileset14ObjectBaseTileID = 200_000*2 + 14*10_000 + 1 // X/tilesets/tileset_14_objects.zel
	tileset15ObjectBaseTileID = 200_000*2 + 15*10_000 + 1 // X/tilesets/tileset_15_objects.zel
	tileset16ObjectBaseTileID = 200_000*2 + 16*10_000 + 1 // X/tilesets/tileset_16_objects.zel
	tileset17ObjectBaseTileID = 200_000*2 + 17*10_000 + 1 // X/tilesets/tileset_17_objects.zel
	// buildings.
	tileset1BuildingBaseTileID  = 200_000*3 + 1*10_000 + 1  // X/tilesets/tileset_1_buildings.zel
	tileset2BuildingBaseTileID  = 200_000*3 + 2*10_000 + 1  // X/tilesets/tileset_2_buildings.zel
	tileset3BuildingBaseTileID  = 200_000*3 + 3*10_000 + 1  // X/tilesets/tileset_3_buildings.zel
	tileset4BuildingBaseTileID  = 200_000*3 + 4*10_000 + 1  // X/tilesets/tileset_4_buildings.zel
	tileset5BuildingBaseTileID  = 200_000*3 + 5*10_000 + 1  // X/tilesets/tileset_5_buildings.zel
	tileset6BuildingBaseTileID  = 200_000*3 + 6*10_000 + 1  // X/tilesets/tileset_6_buildings.zel
	tileset7BuildingBaseTileID  = 200_000*3 + 7*10_000 + 1  // X/tilesets/tileset_7_buildings.zel
	tileset8BuildingBaseTileID  = 200_000*3 + 8*10_000 + 1  // X/tilesets/tileset_8_buildings.zel
	tileset9BuildingBaseTileID  = 200_000*3 + 9*10_000 + 1  // X/tilesets/tileset_9_buildings.zel
	tileset10BuildingBaseTileID = 200_000*3 + 10*10_000 + 1 // X/tilesets/tileset_10_buildings.zel
	tileset11BuildingBaseTileID = 200_000*3 + 11*10_000 + 1 // X/tilesets/tileset_11_buildings.zel
	tileset12BuildingBaseTileID = 200_000*3 + 12*10_000 + 1 // X/tilesets/tileset_12_buildings.zel
	tileset13BuildingBaseTileID = 200_000*3 + 13*10_000 + 1 // X/tilesets/tileset_13_buildings.zel
	tileset14BuildingBaseTileID = 200_000*3 + 14*10_000 + 1 // X/tilesets/tileset_14_buildings.zel
	tileset15BuildingBaseTileID = 200_000*3 + 15*10_000 + 1 // X/tilesets/tileset_15_buildings.zel
	tileset16BuildingBaseTileID = 200_000*3 + 16*10_000 + 1 // X/tilesets/tileset_16_buildings.zel
	tileset17BuildingBaseTileID = 200_000*3 + 17*10_000 + 1 // X/tilesets/tileset_17_buildings.zel
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
		Orientation: "staggered", // isometric(staggered)
		StaggerAxis: "x",
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
	const baseFloorsLayerName = "base_floors"
	baseFloorsTiledTileIDAt := func(x, y int) int {
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
	addLayer(tmxMap, baseFloorsLayerName, baseFloorsTiledTileIDAt)

	// Floors layer.
	const floorsLayerName = "floors"
	floorsTiledTileIDAt := func(x, y int) int {
		floorFrame := int(int16(m.FloorFrameMap[y][x]))
		// tileset type 2 (floors) tile IDs are negative.
		if floorFrame >= 0 {
			return 0 // early return; ignore positive (base floors) and 0.
		}
		floorFrame &= 0x7FFF
		//dbg.Println("negative floorFrame:", floorFrame)
		const base = tileset2FloorBaseTileID // TODO: add support for tileset_NNN/ type.
		tiledTileID := base + floorFrame
		return tiledTileID
	}
	addLayer(tmxMap, floorsLayerName, floorsTiledTileIDAt)

	// Shadows overlay.
	addShadowsOverlay(tmxMap, m)

	// Base walls layer.
	const baseWallsLayerName = "base_walls"
	baseWallFrameAt := make(map[Coordinate]int)
	for _, baseWall := range m.BaseWalls {
		baseWallFrameAt[Coord(int(baseWall.X), int(baseWall.Y))] = int(baseWall.Frame)
	}
	baseWallsTiledTileIDAt := func(x, y int) int {
		baseWallFrame, ok := baseWallFrameAt[Coord(x, y)]
		if !ok {
			return 0
		}
		//dbg.Println("baseWallFrame:", baseWallFrame)
		base := getBaseWallBaseTileID(m)
		tiledTileID := base + baseWallFrame
		return tiledTileID
	}
	addLayer(tmxMap, baseWallsLayerName, baseWallsTiledTileIDAt)

	// Objects layer.
	const objectsLayerName = "objects"
	objectFrameAt := make(map[Coordinate]int)
	for _, object := range m.Objects {
		objectFrameAt[Coord(int(object.X), int(object.Y))] = int(object.Frame)
	}
	objectsTiledTileIDAt := func(x, y int) int {
		objectFrame, ok := objectFrameAt[Coord(x, y)]
		if !ok {
			return 0
		}
		//dbg.Println("objectFrame:", objectFrame)
		const base = tileset2ObjectBaseTileID // TODO: add support for tileset_NNN/ type.
		tiledTileID := base + objectFrame
		return tiledTileID
	}
	addLayer(tmxMap, objectsLayerName, objectsTiledTileIDAt)

	// Buildings layer.
	const buildingsLayerName = "buildings"
	buildingFrameAt := make(map[Coordinate]int)
	for _, building := range m.Buildings {
		buildingFrameAt[Coord(int(building.X), int(building.Y))] = int(building.Frame)
	}
	buildingsTiledTileIDAt := func(x, y int) int {
		buildingFrame, ok := buildingFrameAt[Coord(x, y)]
		if !ok {
			return 0
		}
		//dbg.Println("buildingFrame:", buildingFrame)
		const base = tileset2BuildingBaseTileID // TODO: add support for tileset_NNN/ type.
		tiledTileID := base + buildingFrame
		return tiledTileID
	}
	addLayer(tmxMap, buildingsLayerName, buildingsTiledTileIDAt)
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
	// Layer placed in group to get correct Z-order between layers and overlays
	// (e.g. shadows and mountains).
	group := tmx.Group{
		Name:    tilesLayerName,
		Opacity: 1.0,
		Visible: 1,
	}
	group.Layers = append(group.Layers, tilesLayer)
	tmxMap.Groups = append(tmxMap.Groups, group)
}

// shadowOpacity specifies the opacity of the shadow overlay.
const shadowOpacity = 0.5

// addShadowsOverlay converts MAP shadows overlay to TMX format.
func addShadowsOverlay(tmxMap *tmx.Map, m *maps.Map) {
	// Shadows overlays group.
	const shadowsName = "shadows"
	shadowGroup := tmx.Group{
		Name:    shadowsName,
		Opacity: shadowOpacity,
		Visible: 1,
	}
	for i, shadow := range m.Shadows {
		const tilesetID = 2 // TODO: add support for more tilesets.
		tilesetName := fmt.Sprintf("tileset_%d", tilesetID)
		pngName := fmt.Sprintf("frame_%0004d.png", shadow.Frame)
		pngPath := filepath.Join("..", "tilesets", tilesetName, shadowsName, pngName)
		imageLayer := tmx.ImageLayer{
			Name:    fmt.Sprintf("shadow_%d", i),
			OffsetX: float64(shadow.X),
			OffsetY: float64(shadow.Y),
			Opacity: 1.0,
			Visible: 1,
		}
		img := tmx.Image{
			Source: pngPath,
		}
		imageLayer.Images = append(imageLayer.Images, img)
		shadowGroup.ImageLayers = append(shadowGroup.ImageLayers, imageLayer)
	}
	tmxMap.Groups = append(tmxMap.Groups, shadowGroup)
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

// tilesetInfos specifies the name and dimensions of the following tilesets:
//
//	X/base_floors_tileset.zel                 (base floors)
//	X/base_walls_tileset/base_walls_NNN.zel   (base walls)
//	X/tilesets/tileset_NNN_mountains.zel      (tileset type 0)
//	X/tilesets/tileset_NNN_buildings.zel      (tileset type 1)
//	X/tilesets/tileset_NNN_floors.zel         (tileset type 2)
//	X/tilesets/tileset_NNN_objects.zel        (tileset type 3)
//	X/tilesets/tileset_NNN_shadows.zel        (tileset type 4)
var tilesetInfos = []*TilesetInfo{
	// base floors.
	{
		TilesetName:       "base_floors",
		TilesetWidth:      768,
		TilesetHeight:     2272,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        baseFloorBaseTileID,
	},
	// base walls.
	{
		TilesetName:       "base_walls_1",
		TilesetWidth:      512,
		TilesetHeight:     384,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        baseWalls1BaseTileID,
	},
	{
		TilesetName:       "base_walls_2",
		TilesetWidth:      512,
		TilesetHeight:     384,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        baseWalls2BaseTileID,
	},
	{
		TilesetName:       "base_walls_3",
		TilesetWidth:      512,
		TilesetHeight:     384,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        baseWalls3BaseTileID,
	},
	{
		TilesetName:       "base_walls_4",
		TilesetWidth:      512,
		TilesetHeight:     384,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        baseWalls4BaseTileID,
	},
	{
		TilesetName:       "base_walls_5",
		TilesetWidth:      512,
		TilesetHeight:     384,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        baseWalls5BaseTileID,
	},
	{
		TilesetName:       "base_walls_6",
		TilesetWidth:      512,
		TilesetHeight:     384,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        baseWalls6BaseTileID,
	},
	{
		TilesetName:       "base_walls_7",
		TilesetWidth:      512,
		TilesetHeight:     384,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        baseWalls7BaseTileID,
	},
	// floors (tileset type 2).
	{
		TilesetName:       "tileset_1/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset1FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_2/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset2FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_3/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset3FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_4/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset4FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_5/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset5FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_6/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset6FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_7/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset7FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_8/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset8FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_9/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset9FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_10/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset10FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_11/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset11FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_12/floors",
		TilesetWidth:      768,
		TilesetHeight:     576,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset12FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_13/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset13FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_14/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset14FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_15/floors",
		TilesetWidth:      768,
		TilesetHeight:     800,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset15FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_16/floors",
		TilesetWidth:      768,
		TilesetHeight:     320,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset16FloorBaseTileID,
	},
	{
		TilesetName:       "tileset_17/floors",
		TilesetWidth:      768,
		TilesetHeight:     64,
		TilesetTileWidth:  64, // mapTileWidth
		TilesetTileHeight: 32, // mapTileHeight
		BaseTileID:        tileset17FloorBaseTileID,
	},
	// objects (tileset type 3).
	{
		TilesetName:       "tileset_1/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset1ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_2/objects",
		TilesetWidth:      5376,
		TilesetHeight:     3072,
		TilesetTileWidth:  448,
		TilesetTileHeight: 384,
		BaseTileID:        tileset2ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_3/objects",
		TilesetWidth:      5376,
		TilesetHeight:     3072,
		TilesetTileWidth:  448,
		TilesetTileHeight: 384,
		BaseTileID:        tileset3ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_4/objects",
		TilesetWidth:      5376,
		TilesetHeight:     3072,
		TilesetTileWidth:  448,
		TilesetTileHeight: 384,
		BaseTileID:        tileset4ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_5/objects",
		TilesetWidth:      5376,
		TilesetHeight:     3072,
		TilesetTileWidth:  448,
		TilesetTileHeight: 384,
		BaseTileID:        tileset5ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_6/objects",
		TilesetWidth:      5376,
		TilesetHeight:     3072,
		TilesetTileWidth:  448,
		TilesetTileHeight: 384,
		BaseTileID:        tileset6ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_7/objects",
		TilesetWidth:      5376,
		TilesetHeight:     3072,
		TilesetTileWidth:  448,
		TilesetTileHeight: 384,
		BaseTileID:        tileset7ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_8/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset8ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_9/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset9ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_10/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4576,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset10ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_11/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset11ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_12/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset12ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_13/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset13ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_14/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset14ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_15/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset15ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_16/objects",
		TilesetWidth:      6180,
		TilesetHeight:     4992,
		TilesetTileWidth:  515,
		TilesetTileHeight: 416,
		BaseTileID:        tileset16ObjectBaseTileID,
	},
	{
		TilesetName:       "tileset_17/objects",
		TilesetWidth:      2784,
		TilesetHeight:     496,
		TilesetTileWidth:  232,
		TilesetTileHeight: 248,
		BaseTileID:        tileset17ObjectBaseTileID,
	},
	// buildings (tileset type 1).
	{
		TilesetName:       "tileset_1/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     2560,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset1BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_2/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     4480,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset2BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_3/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     3456,
		TilesetTileWidth:  64,
		TilesetTileHeight: 1152,
		BaseTileID:        tileset3BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_4/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     3456,
		TilesetTileWidth:  64,
		TilesetTileHeight: 1152,
		BaseTileID:        tileset4BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_5/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     2560,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset5BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_6/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     3200,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset6BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_7/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     4800,
		TilesetTileWidth:  64,
		TilesetTileHeight: 800,
		BaseTileID:        tileset7BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_8/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     2816,
		TilesetTileWidth:  64,
		TilesetTileHeight: 704,
		BaseTileID:        tileset8BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_9/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     896,
		TilesetTileWidth:  64,
		TilesetTileHeight: 448,
		BaseTileID:        tileset9BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_10/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     640,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset10BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_11/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     640,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset11BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_12/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     640,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset12BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_13/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     608,
		TilesetTileWidth:  64,
		TilesetTileHeight: 608,
		BaseTileID:        tileset13BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_14/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     896,
		TilesetTileWidth:  64,
		TilesetTileHeight: 448,
		BaseTileID:        tileset14BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_15/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     896,
		TilesetTileWidth:  64,
		TilesetTileHeight: 448,
		BaseTileID:        tileset15BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_16/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     2560,
		TilesetTileWidth:  64,
		TilesetTileHeight: 640,
		BaseTileID:        tileset16BuildingBaseTileID,
	},
	{
		TilesetName:       "tileset_17/buildings",
		TilesetWidth:      6144,
		TilesetHeight:     192,
		TilesetTileWidth:  64,
		TilesetTileHeight: 192,
		BaseTileID:        tileset17BuildingBaseTileID,
	},
}

// addTilesets add all tilesets to the given TMX map.
func addTilesets(tmxMap *tmx.Map) {
	// Add tilesets.
	for _, floorsTilesetInfo := range tilesetInfos {
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
	if tileset.TileWidth != mapTileWidth { // || tileset.TileHeight != mapTileHeight
		// Use tile offset to center tiles larger than 64x32.
		xoff := (mapTileWidth - tileset.TileWidth) / 2
		yoff := 0
		tileset.TileOffset = []tmx.TileOffset{{X: float64(xoff), Y: float64(yoff)}}
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

// Coordinate specifies a (x,y)-coordinate on the map.
type Coordinate struct {
	X int
	Y int
}

// Coord returns the given (x,y)-coordinate.
func Coord(x, y int) Coordinate {
	return Coordinate{
		X: x,
		Y: y,
	}
}

// getBaseWallBaseTileID returns the base walls tileset ID of the given map.
//
//	X/base_walls_tileset/base_walls_NNN.zel
func getBaseWallBaseTileID(m *maps.Map) int {
	// m.BaseWallsTilesetID in range [0, 7)
	switch m.BaseWallsTilesetID {
	case 0:
		return baseWalls1BaseTileID
	case 1:
		return baseWalls2BaseTileID
	case 2:
		return baseWalls3BaseTileID
	case 3:
		return baseWalls4BaseTileID
	case 4:
		return baseWalls5BaseTileID
	case 5:
		return baseWalls6BaseTileID
	case 6:
		return baseWalls7BaseTileID
	default:
		panic(fmt.Errorf("support for base walls tileset ID %d not yet implemented", m.BaseWallsTilesetID))
	}
}
