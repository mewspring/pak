#!/bin/bash

mkdir -p _assets_/tilesets
for i in {1..17}; do
	mkdir -p "_assets_/tilesets/tileset_${i}"
done

# tilset type 1 (buildings)

echo "Creating \"_assets_/tilesets/tileset_1/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_1_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_1/buildings.png

echo "Creating \"_assets_/tilesets/tileset_2/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_2_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_2/buildings.png

echo "Creating \"_assets_/tilesets/tileset_3/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_3_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x1152+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_3/buildings.png

echo "Creating \"_assets_/tilesets/tileset_4/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_4_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x1152+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_4/buildings.png

echo "Creating \"_assets_/tilesets/tileset_5/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_5_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_5/buildings.png

echo "Creating \"_assets_/tilesets/tileset_6/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_6_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_6/buildings.png

echo "Creating \"_assets_/tilesets/tileset_7/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_7_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x800+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_7/buildings.png

echo "Creating \"_assets_/tilesets/tileset_8/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_8_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x704+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_8/buildings.png

echo "Creating \"_assets_/tilesets/tileset_9/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_9_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x448+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_9/buildings.png

echo "Creating \"_assets_/tilesets/tileset_10/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_10_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_10/buildings.png

echo "Creating \"_assets_/tilesets/tileset_11/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_11_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_11/buildings.png

echo "Creating \"_assets_/tilesets/tileset_12/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_12_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_12/buildings.png

echo "Creating \"_assets_/tilesets/tileset_13/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_13_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x608+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_13/buildings.png

echo "Creating \"_assets_/tilesets/tileset_14/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_14_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x448+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_14/buildings.png

echo "Creating \"_assets_/tilesets/tileset_15/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_15_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x448+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_15/buildings.png

echo "Creating \"_assets_/tilesets/tileset_16/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_16_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x640+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_16/buildings.png

echo "Creating \"_assets_/tilesets/tileset_17/buildings.png\""
montage \
	_dump_/X/tilesets/tileset_17_buildings/frame_*.png \
	-tile 96x \
	-background none \
	-gravity south \
	-extent 64x192+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_17/buildings.png
