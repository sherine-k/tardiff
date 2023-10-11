# tardiff
This is a POC utility that generates an archive containing all files modified after a certain date in a selected folder.
It takes 3 parameters as input:
* --starting : a date in format YYYY-MM-dd
* --in: folder that will be scanned for files modified after `starting`
* --out: tar file location for the output

Usage: 
```bash
go build main.go
./main --starting 2023-10-11 --in /home/myuser/myfolder --out diff_10-11.tar
```
The prior command takes all files modified after Oct.11 in /home/myuser/myfolder and generating diff_10-11.tar from those files.
