# PhotoPainter (B) Image Converter

A simple command-line tool to convert images for use with PhotoPainter (B). This tool converts common image formats to 24-bit BMP format with 6 colors (black, white, red, green, blue, yellow) compatible with PhotoPainter (B).

## Features

- Supports common image formats (JPG, PNG, GIF, BMP, etc.)
- Floyd-Steinberg dithering algorithm for color conversion
- Automatic resizing to 800×480 resolution
- Automatic rotation of portrait images (90-degree rotation for proper orientation)
- Batch processing of multiple files in a directory
- Sequential file numbering in batch mode (0001_0001.bmp, 0001_0002.bmp, etc.)

## Usage

```
PhotoPainter (B) Image Converter
Usage: photoconvert [options] input_file/directory

Options:
  -o <dir>       Output directory (required for batch mode)
  -batch         Batch mode (process all images in a directory)
  -depth <n>     Maximum subdirectory exploration depth for batch mode (default: 3)
  -r <res>       Resolution (800x480 or 480x800, default: 800x480)
  -rotate=false  Disable automatic rotation of portrait images (default: enabled)
  -v             Verbose logging
  -h             Display this help message

Examples:
  Single file conversion:
    photoconvert input.jpg
  Specify output directory:
    photoconvert -o /path/to/output input.jpg
  Batch processing:
    photoconvert -batch -o /path/to/output /path/to/images/
  Batch processing including subdirectories:
    photoconvert -batch -depth 3 -o /path/to/output /path/to/images/
  All options:
    photoconvert -batch -o /path/to/output -depth 2 -r 800x480 -v /path/to/images/
```

## Batch Processing Behavior

When using batch mode:
1. All files in each directory are sorted alphabetically before processing
2. Directories are assigned sequential numbers (0001, 0002, etc.)
3. Files within each directory are numbered sequentially (0001, 0002, etc.)
4. All output files are written to a single flat directory (specified by -o)
5. Output filenames use the format: `DDDD_FFFF.bmp` where:
   - `DDDD` is the directory sequence number
   - `FFFF` is the file sequence number within that directory

## Installation

Clone the repository and build:

```bash
git clone https://github.com/fuba/convert2photopainter_b.git
cd convert2photopainter_b
go build -o photoconvert ./cmd/convert2photopainter/
```

## Project Structure

```
/convert2photopainter_b/
  ├── cmd/
  │   └── convert2photopainter/
  │       └── main.go            # Command-line interface
  ├── internal/
  │   ├── convert/
  │   │   └── convert.go         # Image conversion processing
  │   ├── dither/
  │   │   └── dither.go          # Floyd-Steinberg dithering
  │   └── resize/
  │       └── resize.go          # Image resizing processing
  ├── go.mod                     # Dependency management
  └── README.md                  # Project explanation
```

## Process Flow

1. Load the image
2. Resize to the specified resolution (default: 800×480)
3. Apply Floyd-Steinberg dithering algorithm to reduce to 6 colors
4. Save as 24-bit BMP format

## Notes

- The output image color count is limited to 6 colors (black, white, red, green, blue, yellow) according to PhotoPainter (B) specifications
- Floyd-Steinberg dithering is used for color conversion to achieve a more natural appearance
- Images are automatically resized, but may be cropped if the aspect ratio differs from the target
- When creating a TF card on a MAC, hidden files may be generated, so it is recommended to delete fileList.txt and index.txt from the root directory of the TF card after conversion, as well as hidden files in the pic folder

## References

- [PhotoPainter (B) Official Wiki](https://www.waveshare.com/wiki/PhotoPainter_(B))