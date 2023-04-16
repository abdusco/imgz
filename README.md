# imgz

A utility to resample images.

## Requirements
- `vips`  
  This is a soft requirement. Instead of linking `libvips`,  
  this app uses [`vips`](https://github.com/libvips/libvips) CLI app to resize images.
  If it's not available in `$PATH`, it will default to [`imaging`](https://github.com/disintegration/imaging) — a pure image processing library written in golang.
  It uses unfortunately a lot more RAM and CPU than vips.

## Usage

```shell
Usage: imgz <command>

Flags:
  -h, --help       Show context-sensitive help.
      --debug      Enable debug logging
      --version    Show version and exit

Commands:
  resize <image-path>
    Resize an image

  resize-dir --output-dir=STRING <source-dirs> ...
    Resize a folder of images
```

### Resize a single image

```shell
Usage: imgz resize <image-path>

Resize an image

Arguments:
  <image-path>    Path to image. Use "-" for stdin

Flags:
  -h, --help             Show context-sensitive help.
      --debug            Enable debug logging
      --version          Show version and exit

  -o, --output=STRING    Path to output file. Use "-" for stdout. Defaults to $sourceDir/resized/$source.jpg or stdout if input is stdin
      --max-size=5000
      --quality=75
```

Example:

```shell
imgz resize --max-length 5000 --quality 80 img.jpg -o out.jpg
```

### Resize whole dirs containing images

```shell
Usage: imgz resize-dir --output-dir=STRING <source-dirs> ...

Resize a folder of images

Arguments:
  <source-dirs> ...    List of paths to image folders

Flags:
  -h, --help                 Show context-sensitive help.
      --debug                Enable debug logging
      --version              Show version and exit

  -o, --output-dir=STRING    Dir to save zip files
      --max-size=5000        Max side length of resized images
      --quality=75           JPEG quality
      --clean                Delete source dirs after resizing
```

```shell
imgz resize-dir --max-length 5000 --quality 80 --output-dir /output --clean ./sourcefolder1 ./sourcefolder2 ./sourcefolder3
```