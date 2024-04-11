<p align="center">
<img src="src/images/footprint.png" alt="drawing" width="250" height="250" />
</p>

#  File Footprint Checker
File Footprint Checker is a software to compute recursively files footprints (md5sum) of all files in a list of directories

# Quick start
- edit the config/path.csv file
- enter one directory path to scan, per line (no header in this table)
- start the software

# Key characteristics
- unlimited number of directory path
- md5 sum calculation
- TSV output with 2 columns signatures and path
- statically compiled (written in Go), nothing to install 

# ScreenShots
![CLI](src/images/screenshot.png)

Results (md5 - path):

826e8142e6baabe8af779f5f490cf5f5	test/file1.txt

1c1c96fd2cf8330db0bfa936ce82f3b9	test2/file2.txt