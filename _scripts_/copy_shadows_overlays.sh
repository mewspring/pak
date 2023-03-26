#!/bin/bash

mkdir -p _assets_/tilesets
for i in {1..17}; do
	mkdir -p "_assets_/tilesets/tileset_${i}/shadows"
done

# tilset type 4 (shadows)

for i in {1..17}; do
	echo "Creating \"_assets_/tilesets/tileset_${i}/shadows/\""
	cp _dump_/X/tilesets/tileset_${i}_shadows/frame_*.png _assets_/tilesets/tileset_${i}/shadows/
done
