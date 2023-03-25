#!/bin/bash

mkdir -p _assets_/tilesets

for i in {1..17}; do
	mkdir -p "_assets_/tilesets/tileset_${i}"
	# tilset type 0 (mountains and stairs)
	# tilset type 1 (walls and buildings)
	# tilset type 2 (floors)
	echo "Creating \"_assets_/tilesets/tileset_${i}/floors.png\""
	montage \
		_dump_/X/tilesets/tileset_${i}_floors/frame_*.png \
		-tile 12x \
		-background none \
		-gravity south \
		-extent 64x32+0+0 \
		-geometry +0+0 \
		_assets_/tilesets/tileset_${i}/floors.png
	# tilset type 3 (objects)
	# tilset type 4 (shadows)
done

echo 'Creating "_assets_/tilesets/base_floors.png"'
montage \
	_dump_/X/base_floor_tileset/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 64x32+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/base_floors.png
