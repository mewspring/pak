#!/bin/bash

mkdir -p _assets_/tilesets
for i in {1..17}; do
	mkdir -p "_assets_/tilesets/tileset_${i}/backgrounds"
done

# tilset type 0 (backgrounds)

for i in {1..17}; do
	# skip empty files (X/tilesets/tileset_{3,4,5,6}_backgrounds.zel)
	if [[ $i -eq 3 ]]; then
		continue
	elif [[ $i -eq 4 ]]; then
		continue
	elif [[ $i -eq 5 ]]; then
		continue
	elif [[ $i -eq 6 ]]; then
		continue
	fi
	echo "Creating \"_assets_/tilesets/tileset_${i}/backgrounds/\""
	cp _dump_/X/tilesets/tileset_${i}_backgrounds/frame_*.png _assets_/tilesets/tileset_${i}/backgrounds/
done
