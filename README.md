# img-sort

Introducing `img-sort`, a little golang tool that helps you keep your photos and videos organized. This tool will sort all your media in year and month folders, making it easy for you to find the photos and videos you need.

## Features

- Simple and easy to use with flags
- Supports various file formats including `JPEG`, `HEIC`, `MP4`, `MOV` and more
- Reads the exif data if available otherwise uses the modification time
- Photos and videos without date/time information are moved to a directory `unknown`
- Duplicates are handled by appending a postfix `-1`, `-2`, `-3` and so on

## Folder Structure

The folder structure created by `img-sort` will look like the following example:

```yaml
├── 2021
│   ├── 2021-01
│   │   ├── 2021-01-07_11.23.44.jpg
│   │   ├── 2021-01-07_13.24.53.heic
│   │   ├── 2021-01-07_20.27.47.jpg
│   │   ├── 2021-01-07_20.27.47-1.jpg
│   │   ├── 2021-01-09_15.58.24.jpg
│   │   ├── 2021-01-09_21.39.27.mp4
│   │   ├── 2021-01-09_23.13.37.jpg
│   │   ├── 2021-01-14_10.44.50.mov
│   ├── 2021-02
│   |   ├─ ...
│   ├── ...
│   └── 2021-12
├── 2022
│   ├── 2022-01
│   ├── 2022-02
│   ├── ...
│   └── 2022-12
├── unknown
├── ...
```

## Installation

In order to use this application, the following dependencies must be installed on your system:

- Perl: A version of Perl `>= 5`and `< 6` must be installed.
- Exiftool: A version of Exiftool `>= 12.55` must also be installed on your system. This application searches for the binaries at `/usr/bin/exiftool`.

Please ensure that these dependencies are installed before attempting to run this application.

### Windows

```bash
choco install perl --version 5.30
choco install exiftool --version 12.55
```

### MacOS

```bash
brew install perl@5.30
brew install exiftool@12.55
```

### Linux

```bash
sudo apt-get install perl=5.30
sudo apt-get install exiftool=12.55
```

After all dependencies are installed the executable of `img-sort` can be added to the `$PATH` variable or moved from `build/img-sort` to `/usr/bin/img-sort`.

## Usage

Simply run `img-sort` and set `--source` and `--target`. The tool will take care of the rest.

```bash
img-sort --source /path/to/source --target /path/to/target
```

If there is no date information available for given photos or videos the file modification time can be used as a fallback.

```bash
img-sort --source /path/to/source --target /path/to/target --modtime
```

### Available flags

- `--version`: Version info
- `--source`: Source path
- `--target`: Target path
- `--modtime`: Modification time fallback

## Building from sources

1. Clone this repository: `git clone https://github.com/patrickap/img-sort.git`
2. Navigate into the project directory: `cd img-sort`
3. Build the tool: `sh build.sh`
4. Run the binary: `build/img-sort`
