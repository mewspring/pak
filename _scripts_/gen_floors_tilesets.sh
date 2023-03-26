#!/bin/bash

mkdir -p _assets_/tilesets
for i in {1..17}; do
	mkdir -p "_assets_/tilesets/tileset_${i}"
done

# base floors tileset

echo 'Creating "_assets_/tilesets/base_floors.png"'
montage \
	_dump_/X/base_floors_tileset/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 64x32+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/base_floors.png

# tilset type 2 (floors)

for i in {1..17}; do
	mkdir -p "_assets_/tilesets/tileset_${i}"
	echo "Creating \"_assets_/tilesets/tileset_${i}/floors.png\""
	montage \
		_dump_/X/tilesets/tileset_${i}_floors/frame_*.png \
		-tile 12x \
		-background none \
		-gravity south \
		-extent 64x32+0+0 \
		-geometry +0+0 \
		_assets_/tilesets/tileset_${i}/floors.png
done
