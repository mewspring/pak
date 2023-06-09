# pak

## Install

```bash
git clone https://github.com/mewspring/pak
cd pak
go install ./cmd/pak_dump
go install ./cmd/zel_patch
go install ./cmd/zel_dump
go install ./cmd/map_dump
```

## Usage

```bash
# Extract PAK archive.
pak_dump -listfile listfile.json X.PAK
```

```bash
# Patch broken ZEL images.
zel_patch
```

```bash
# Convert ZEL images to PNG format.
find ./_dump_/X -type f -name "*.zel" -exec zel_dump -pal _dump_/X/core/core.pal {} \;
```

```bash
# Generate tileset sprite sheets.
./_scripts_/gen_tilesets.sh
```

```bash
# Convert MAP files to TMX format.
map_dump _dump_/X/tilesets/map_*.map
```
