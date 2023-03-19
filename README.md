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
pak_dump -listfile listfile.json X.PAK
```

```bash
zel_dump -pal _dump_/X/core/core.pal _dump_/X/path/to/image.zel
```

```bash
# Convert ZEL image to PNG format.
#
# NOTE: _dump_/X/archive_0000/file_0002.bin (file size 1024) is an RGBA palette
# with 256 colours.
#
# NOTE: _dump_/X/archive_0002/archive_0035.bin is a ZEL image.
zel_dump -pal _dump_/X/archive_0000/file_0002.bin _dump_/X/archive_0002/archive_0035.bin
```

*Note*: since the header format of PAK files and ZEL images is identical, these two file formats are indistinguishable by header data.

As such, `pak_dump` extracts the frame contents of `X/archive_0002/archive_0035.bin` to `X/archive_0002/archive_0035/file_0000.bin` (even though `archive_0035.bin` is a ZEL image and not a PAK archive).

Since, `zel_dump` requires the ZEL frame offset table to parse ZEL images, we must use `X/archive_0002/archive_0035.bin` and not `X/archive_0002/archive_0035/file_0000.bin`.
