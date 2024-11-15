# mango peeler

... is a CLI tool intended for locating and removing images inserted by aggregator sites from your local storage of (scraped) manga. The list of aggregator images the script will be looking for (grouped by site), can be found [here](https://github.com/valdotle/mangopeeler/tree/main/images). Furthermore, it's capable of detecting and removing duplicate images in the same directory.

I've also been thinking of various other features that I might add at a later point, such as:

-   [ ] deleting empty directories
-   [ ] renaming images (of a directory) to have incrementing numbers starting at 1, to patch gaps resulting from deleting images (or other reasons)
-   [ ] deleting images that can't be read and thus are most likely corrupted/malformed
-   [ ] filter (and remove) files by filename/matching patterns and regex/file extension/...
-   [ ] support to filter for custom images not part of the script

# usage

To get started, download the script for your corresponding OS [here](https://github.com/valdotle/mangopeeler/releases). While the script _should_ be cross-platform compatible, I'm only able to test on Windows, so keep that in mind.

Currently, **unzipping** folders **isn't supported**. That means, the script will only be able to process **images** that are simply **located in a directory**. When coding, I was assuming a directory structure like so:

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

See [the next section](#performance-considerations) for a slightly more detailed breakdown of certain, performance-relevant, configuration options.

### command flags

If you want more control over the script, it supports a variety of command flags, all of which can be viewed with the `-h`/`-help` flag.

-   #### `-del`/`-delete`
    Whether you want to delete duplicate and aggregator images (default `true`)
-   #### `-dt`/`-directory-threads`
    How many directory entries to process simultaneously (default `20`), only applicable if walking subdirectories is enabled. Set to `1`/`0` to disable.
-   #### `-dir`
    The directory to execute this script in (defaults to the current directory)
-   #### `-det`/`-directory-entry-threads`
    How many directory entries/files to process simultaneously (default `5`). Set to `1`/`0` to disable.
-   #### `-l`/`-log`
    Whether to create a logfile for actions performed by the script (default `true`)
-   #### `-lat`/`-log-at`
    Where to store logfiles, if logging is enabled. Defaults to a folder named `mango peels` in the current directory.
-   #### `-s`/`-site`
    Which aggreagtor site(s)'s images to search for. You can set this flag multiple times to supply multiple values. By default, all aggregator sites' images will be used.
-   #### `-w`/`-walk`
    Whether to recursively process subdirectories, if there are any (default `true`)

### config file

All options can be set using a config file instead as well. The `config.json` file **must be located in the same directory as the script** to take effect. Command flags will take precedence over the config file, meaning you can use it to store your baseline settings and set command flags to overrule them on the fly as needed. A sample `config.json` file, set to the [command flag](#command-flags) defaults comes with all binaries and can be found [here](https://github.com/valdotle/mangopeeler/tree/main/config.json) as well.

# performance considerations

While the configuration allows for quite some tweaking performance wise, the current defaults are based purely on my educated guesses and limited testing. I should probably do some (proper) benchmarking at some point to find good defaults. That being said, the exact tuning heavily depends on your hardware specifications (primarily disk reading speed), the length of chapters, the average image size, the overall directory structure you want to scan and how much resources you want the script to use, to name just a few things.

Here are some theoretical performance considerations when it comes to `directory-entry-threads`, `directory-threads`, `site` and `walk` though:

<details>
  <summary><h3> `site` and `walk`</h3></summary>
  If processing subdirectories is enabled with `walk` (the default), the script will not only scan the current/specified directory, but all its subdirectories as well. That means more directory entries to look through and thus naturely more compute. Make sure to set `walk` to `false` if you only mean to scan a specific directory's contents to avoid wasting resources.

Similarly, `site` allows to limit which aggregator site's images to look for. By default, the script will search all aggregator images. But if you are going to scan a bunch of Vietnamese manga, a Portugese aggregator's images are rather unlikely to appear - and vice versa. By specifying which aggregator's images to look for, you can reduce the time it takes to match an image against the aggregator images.

</details>

<details>
  <summary><h3> `directory-threads` and `directory-entry-threads`</h3></summary>
  Directory entry threads is the number of threads allocated per directory. That is, if the folder for a chapter contains X files, each directory entry thread can process one of them simultaneously. Increasing the number of directory entry threads might be useful when processing directories with many entries (think long chapters). On the other hand, titles with short chapters (=few entries per directory) might be unable to utilize multiple threads processing directory entries. In this case, reducing directory entry threads or even disabling threading on directory level might be beneficial. Additionally, managing threads introduces a certain performance overhead which may not outweigh the performance gains earned from threading so using a very low number of directory entry thread will probably your reduce performance as well.

Similarly, directory threads determine the number of directories (think chapters) that can be processed simultaneously. While increasing directory threads is mainly a means to increase resource utilization, decreasing or disabling directory entry threads is useful when processing a single/small number of directories only or to run the script in the background.

The total number of threads is simply the product of `directory-threads` and `directory-entry-threads`. This number is the main performance impacting metric. Not only will it increase CPU load, and memory usage, more threads will also increase the number of files read at the same time resulting in higher disk utilization. I'd suggest you find a value that either meets your resource usage limits or results in 100% disk usage using the task manager (or similar tooling). In this case, the maximum read speed of your (hard)drive becomes the limitting factor to the script's performance, so adding more threads will mainly introduce overhead instead of increasing productivity.

For a given number of total threads, the ratio of `directory-threads` and `directory-entry-threads` should mainly have an effect on memory usage, if checking for duplicates within a directory is enabled:<br>
Increasing the number of directory threads while reducing the number of directory entry threads will lead to an increase in memory usage and vice versa.

</details>

### TLDR

-   Disabling stuff reduces resources used.
-   More threads = higher resource utilization. Don't increase threads further once disk usage hits 100%.
-   Longer chapters → consider increasing the number of directory entry threads.
-   For a given total number of threads, allocating more threads to `directory-threads` while decreasing `directory-entry-threads` should decrease memory usage.

The inverse is true for all of the above.

Very low values for either `directory-threads` or `directory-entry-threads` can hamper performance rather than improving it. Consider disabling the respective option.

# notes

While I've aimed for this script to run as stable as possible (or did I?), I won't make any guarantees. Feel free to open an [issue](https://github.com/valdotle/mangopeeler/issues) if you run into any problems though!

# contributing

Contributions are generally welcome. Adding widespread aggregator images would be the potentially simplest and most beneficial contribution to mention here.

If you're intending to make larger changes/additions to the overall codebase, it'd be desirable to open a [PR](https://github.com/valdotle/mangopeeler/pulls)/issue (or any other medium for conversation) beforehand, so we can avoid time being invested into changes that ultimately might not get merged.
