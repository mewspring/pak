# pak

## Install

```bash
git clone https://github.com/mewspring/pak
cd pak
go install ./cmd/pak_dump
go install ./cmd/zel_dump
```

## Usage

```bash
# Extract PAK archive.
pak_dump X.PAK
```

```bash
# Convert ZEL image to PNG format.
zel_dump -pal _dump_/X/archive_0000/file_0002.bin _dump_/X/archive_0002/archive_0035.bin
```
