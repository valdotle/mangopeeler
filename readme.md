# mango peeler

... is a CLI tool intended for locating and removing images inserted by aggregator sites from your local storage of (scraped) manga. The list of aggregator images the script will be looking for (grouped by site), can be found [here](https://github.com/valdotle/mangopeeler/tree/main/images). Furthermore, it's capable of detecting and removing duplicate images in the same directory.

I've also been thinking of various other features that I might add at a later point, such as:

-   [x] deleting empty directories
-   [ ] renaming images (of a directory) to have incrementing numbers starting at 1, to patch gaps resulting from deleting images (or other reasons)
-   [ ] deleting images that can't be read and thus are most likely corrupted/malformed
-   [ ] filter (and remove) files by filename/matching patterns and regex/file extension/...
-   [x] support to filter for custom images not part of the script

# usage

To get started, download the script for your corresponding OS [here](https://github.com/valdotle/mangopeeler/releases/latest). While the script _should_ be cross-platform compatible, I'm only able to test on Windows, so keep that in mind.

Currently, **unzipping** folders **isn't supported**. That means the script will only be able to process **images** that are simply **located in a directory**. When coding, I was assuming a directory structure like so:

```
manga
|
|_ title a
|   |_ chapter 1
|   |   |_ page 1.img
|   |   |_ page 2.img
|   |   ...
|   |
|   |_ chapter 2
|   ...
|
|_ tile d
...
```

where the "manga" directory would be the entry point for your script. That being said, you could start at the title level as well, or even run the script on a specific chapter only. Even having all images in a single directory would work, if chapters are compiled by volume for example.

Once downloaded, you can simply move the script into a directory where you want it to remove aggregator and duplicate images. When executed, it will start searching and removing any matching images from the directory **and its subdirectories**.

Note that currently **only `.png`, `.gif`, `.jpg` and `.jpeg` images are supported**. This doesn't mean that you can't run the script in directories that contain any other types of images or anything. Just that **those images will be skipped** and remain unaffected. I might expand the functionality in that regard (looking at `.webp` especially) in the future though, if it turns out to be necessary.

# configuration

You can configure the behaviour of the script using either command flags, a config file or both. An explanation of the individual settings available can be found below.

## command flags

If you want more control over the script, it supports a variety of command flags, all of which can be viewed with the `-h`/`-help` flag.

-   ### `-c`/`-custom-images`
    Provide the path to a directory containing custom images to check for, similar to the aggregator ones included in the script. Defaults to a folder named `custom` in the current directory.
-   ### `-del`/`-delete`
    Whether you want to delete directory entries matching your searching criteria (default `true`)
-   ### `-dir`/`-directory`
    The directory to execute this script in (defaults to the current directory)
-   ### `-dup`/`-duplicates`
    Whether to check for duplicate images within a directory (default `false`). Not enabled by default, because there might be [false positives](#false-positives).<br>
    If there are two identical images within the same directory, the more inset (closer to the middle of the directory) one will be the one flagged as duplicate (the intention behind this was to remove duplicate credit pages in the middle of a chapter).
-   ### `-edr`/`-empty-directory`
    Whether to check for empty directories (default `true`)
-   ### `-l`/`-log`
    Whether to create a logfile for actions performed by the script (default `true`)
-   ### `-lat`/`-log-at`
    Where to store logfiles, if logging is enabled. Defaults to a folder named `mango peels` in the current directory.
-   ### `-s`/`-site`
    Which aggregator site(s)'s images to search for. You can set this flag multiple times to supply multiple values. By default, all aggregator sites' images will be used. Set to `none` to disable.
-   ### `-t`/`-threads`
    How many directories to process simultaneously (default `10`). Only applicable if walking subdirectories is enabled. Set to `1`/`0` to disable.
-   ### `-w`/`-walk`
    Whether to recursively process subdirectories, if there are any (default `true`)

## config file

All options can be set using a config file instead as well. The `config.json` file **must be located in the same directory as the script** to take effect. Command flags will take precedence over the config file, meaning you can use it to store your baseline settings and set command flags to overrule them on the fly as needed. A sample `config.json` file, reflecting the [command flag](#command-flags) defaults comes with all binaries and can be found [here](https://github.com/valdotle/mangopeeler/tree/main/config.json) as well.

# false positives

The script uses a similarity metric to determine images matching aggregator images and duplicate images. Since exact byte equality is pretty useless, I'm resorting to this method which in turn means, there's a chance for images to be flagged incorrectly. Especially the duplicate check is prone to mistakes when there are images with few details (prone meaning I've had one case of an almost [completely black image with a different word in the middle so far](https://github.com/valdotle/mangopeeler/tree/main/false%20positives)). I'll have to see if/to which extent I can optimize the similarity metric to minimize the number of those false positives. In the meantime, you can disable deleting and double check the duplicates found by the script with the log output first, if you want to be absolutely safe.

# performance

Adjusting the number of threads has the most impact on performance. While more threads allow to process more directories simultaneously, this comes with increased resource usage as well.

When trying to maximize the file processing speed, the limiting factor will usually be your disk's reading speed. Use the task manager (or similar tooling) to check when disk utilization reaches 100%. At this point, increasing the number of threads will become less effective or even counter productive. Obviously, you can also run the script with a lower number of threads to consume less resources, if you only want it to run in the background for example.

The current number of threads is simply an educated guess made by me based on my limited testing. I _should_ probably do some benchmarking to find an optimal default at some point. Then again, there are many factors that are very dependent on your hardware specs and optimization goals so those probably would have limited use anyway.

# issues

While I've aimed for this script to run as stable as possible (or did I?), I won't make any guarantees. Feel free to open an [issue](https://github.com/valdotle/mangopeeler/issues) if you run into any problems though!

# contributing

Contributions are generally welcome. Adding widespread aggregator images would be the potentially simplest and most beneficial contribution to mention here.

If you're intending to make larger changes/additions to the overall codebase, it'd be desirable to open a [PR](https://github.com/valdotle/mangopeeler/pulls)/issue (or any other medium for conversation) beforehand, so we can avoid time being invested into changes that ultimately might not get merged.
