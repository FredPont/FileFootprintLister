<p align="center">
<img src="src/images/footprint.png" alt="drawing" width="250" height="250" />
</p>

#  File Footprint Lister
File Footprint Lister is a software to compute recursively files footprints (md5sum, sha256, xxhash64, murmmur64, cityhash64, cityhash128, clickhouse64, clickhouse128) of all files in a list of directories.

Lists can be compared using [CompareFootprintLists](https://github.com/FredPont/CompareFootprintLists)

# Quick start
- edit the config/path.csv file
- enter one directory path to scan, per line (no header in this table)
- to exclude some files or directories, edit the config/exclude.csv file : one string to exclude in paths per line.
If this string is found in the path, this path will be skipped. [.git ~snapshot] are excluded by default.


If the path contains spaces, commas... it can be necessary to quote the path : "my,path with, spaces" 
- start the software in the FileFootprintLister directory
```
Usage :

  -a string

        algorithm to use. md5, xxhash, murmmur, cityhash64, cityhash128, clickhouse64, clickhouse128 or sha256 (default "md5")

  -n int

    	number of CPUs for parallel file processing (default 8).
      The optimal number of CPU depends on the CPU and disk speed. 
      Top speed is generally obtained with 8-16 CPUs [see Benchmarks]

example :

./ffpl-x86_64_linux.bin                   # md5 sum computation with 8 files processed in parallel by default

./ffpl-x86_64_linux.bin -n 14             # md5 sum computation with 14 files processed in parallel

./ffpl-x86_64_linux.bin -a md5            # md5 sum computation

./ffpl-x86_64_linux.bin -a sha256         # sha256 sum computation

./ffpl-x86_64_linux.bin -a xxhash         # xxhash64 sum computation

./ffpl-x86_64_linux.bin -a murmur         # murmur64 sum computation

./ffpl-x86_64_linux.bin -a cityHash64     # cityhash64 sum computation

./ffpl-x86_64_linux.bin -a clickhouse64   # clickhouse64 sum computation
```
- the result tables in TSV are in the result directory. The output table has 3 columns : 
  - footprint
  - file name
  - file path

# Key characteristics
- unlimited number of directory path
- parallel file processing
- md5 sum calculation
- sha256 sum calculation
- murmur64 sum calculation
- xxhash64 sum calculation
- TSV output with 3 columns signatures, name and path
- statically compiled (written in Go), nothing to install 

# Benchmarks (v20241113)
![CLI](src/benchmark/benchmark.png)

Speed : 7.4 files/sec - 9.2 GB/min

![ALGO](src/benchmark/boxplot.png)

# ScreenShots

![CLI](src/images/screenshot.png)

!
