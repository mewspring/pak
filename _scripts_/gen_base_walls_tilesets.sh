#!/bin/bash

mkdir -p _assets_/tilesets

# base walls tilesets

for i in {1..7}; do
	echo "Creating \"_assets_/tilesets/base_walls_${i}.png\""
	montage \
		_dump_/X/base_walls_tileset/base_walls_${i}/frame_*.png \
		-tile 8x \
		-background none \
		-gravity south \
		-extent 64x192+0+0 \
		-geometry +0+0 \
		_assets_/tilesets/base_walls_${i}.png
done
