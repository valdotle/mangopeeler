# mango peeler

... is a CLI tool intended for locating and removing aggregator images from your local storage of (scraped) manga. The list of aggregator images the skript will be looking for (sorted by site), can be found [here](https://github.com/valdotle/mangopeeler/tree/main/images). Furthermore, it's capable of detecting duplicate files in the same directory.

I've also been thinking of various other features that I might add at a later point, such as:

-   [ ] deleting empty directories
-   [ ] renaming images (of a directory) to have incrementing numbers starting at 1, to patch gaps resulting from deleting images (or other reasons)
-   [ ] deleting images that can't be read and thus are most likely corrupted/malformed
-   [ ] filter (and remove) files by filename/matching patterns and regex/file extension/...
-   [ ] support to filter for custom images not part of the skript

# usage

To get started, download the skript for your corresponding OS [here](https://github.com/valdotle/mangopeeler/releases). While the skript _should_ be cross-platform compatible, I'm only able to test on Windows, so keep that in mind.

Once downloaded you can simply move the skript into a directory you want it to check for duplicates and it will automatically search and remove any matching images from the directory **and its subdirectories**.

# configuration

## command flags

If you want more control over the skript, it supports a variety of command flags, all of which can be viewed setting the `-h`/`-help` flag.

-   ### `-d`/`-delete`
    Whether you want to delete duplicate files (default `true`)
-   ### `-dt`/`-directory-entry-threads`
    How many directory entries to process simultaneously (default 10), only applicable if walking subdirectories is enabled
-   ### `-dir`
    The directory to execute this script in (defaults to the current directory)
-   ### `-det`/`-directory-entry-threads`
    How many directory entries to process simultaneously (default 10)
-   ### `-l`/`-log`
    Whether to create logfiles for actions performed by the script (default true)
-   ### `-lat`/`-log-at`
    Where to store logfiles, if enabled; defaults to a folder named "mango peels" in the current directory.
-   ### `-s`/`-site`
    Which site(s)'s images to check for duplicates. You can set this flag multiple times to supply multiple values. By default, all aggregator images will be filtered for.
-   ### `-w`/`-walk`
    Whether to walk subdirectories, if there are any (default true)

## config file

All options can be set using a config file instead as well. Command flags will take precedence over the config file, meaning you can use the config file to store your baseline settings and set command flags to overrule them on the fly as needed. A sample config file, mirroring the [command flag](#command-flags) defaults can be found [here](https://github.com/valdotle/mangopeeler/tree/main/config.json).
