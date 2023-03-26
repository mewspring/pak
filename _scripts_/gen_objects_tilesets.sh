#!/bin/bash

mkdir -p _assets_/tilesets
for i in {1..17}; do
	mkdir -p "_assets_/tilesets/tileset_${i}"
done

# tilset type 3 (objects)

echo "Creating \"_assets_/tilesets/tileset_1/objects.png\""
montage \
	_dump_/X/tilesets/tileset_1_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_1/objects.png

echo "Creating \"_assets_/tilesets/tileset_2/objects.png\""
montage \
	_dump_/X/tilesets/tileset_2_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 448x384+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_2/objects.png

echo "Creating \"_assets_/tilesets/tileset_3/objects.png\""
montage \
	_dump_/X/tilesets/tileset_3_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 448x384+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_3/objects.png

echo "Creating \"_assets_/tilesets/tileset_4/objects.png\""
montage \
	_dump_/X/tilesets/tileset_4_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 448x384+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_4/objects.png

echo "Creating \"_assets_/tilesets/tileset_5/objects.png\""
montage \
	_dump_/X/tilesets/tileset_5_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 448x384+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_5/objects.png

echo "Creating \"_assets_/tilesets/tileset_6/objects.png\""
montage \
	_dump_/X/tilesets/tileset_6_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 448x384+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_6/objects.png

echo "Creating \"_assets_/tilesets/tileset_7/objects.png\""
montage \
	_dump_/X/tilesets/tileset_7_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 448x384+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_7/objects.png

echo "Creating \"_assets_/tilesets/tileset_8/objects.png\""
montage \
	_dump_/X/tilesets/tileset_8_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_8/objects.png

echo "Creating \"_assets_/tilesets/tileset_9/objects.png\""
montage \
	_dump_/X/tilesets/tileset_9_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_9/objects.png

echo "Creating \"_assets_/tilesets/tileset_10/objects.png\""
montage \
	_dump_/X/tilesets/tileset_10_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_10/objects.png

echo "Creating \"_assets_/tilesets/tileset_11/objects.png\""
montage \
	_dump_/X/tilesets/tileset_11_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_11/objects.png

echo "Creating \"_assets_/tilesets/tileset_12/objects.png\""
montage \
	_dump_/X/tilesets/tileset_12_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_12/objects.png

echo "Creating \"_assets_/tilesets/tileset_13/objects.png\""
montage \
	_dump_/X/tilesets/tileset_13_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_13/objects.png

echo "Creating \"_assets_/tilesets/tileset_14/objects.png\""
montage \
	_dump_/X/tilesets/tileset_14_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_14/objects.png

echo "Creating \"_assets_/tilesets/tileset_15/objects.png\""
montage \
	_dump_/X/tilesets/tileset_15_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_15/objects.png

echo "Creating \"_assets_/tilesets/tileset_16/objects.png\""
montage \
	_dump_/X/tilesets/tileset_16_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 515x416+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_16/objects.png

echo "Creating \"_assets_/tilesets/tileset_17/objects.png\""
montage \
	_dump_/X/tilesets/tileset_17_objects/frame_*.png \
	-tile 12x \
	-background none \
	-gravity south \
	-extent 232x248+0+0 \
	-geometry +0+0 \
	_assets_/tilesets/tileset_17/objects.png
