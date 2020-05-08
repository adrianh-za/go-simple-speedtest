# go-simple-speedtest
Program written in Go to perform simple file download to calculate download speed.

When the program is executed, it appends the download data to a file (speedtest.csv) in CSV format.

The CSV column data is:
* time of file download
* total bytes downloaded
* total megabytes downloaded
* total milliseconds to complete download
* total seconds to complete download
* speed of download in megabits/second (Mbps)
* speed of download in megabytes/second (MBps)

## Usage ##

1) go get github.com/adrianh-za/go-simple-speedtest
2) browse to $/go/src/github.com/adrianh-za/go-simple-speedtest
3) go run main.go

## Extra ##

The following cron job runs every 4 hours.  The cron job has multiple steps:

* <b>0 */4 * * *</b> = Run every 4 hours.
* <b>PATH=$PATH:/usr/local/go/bin</b> = Set the path to include GO executable location.
* <b>cd /home/pi/Documents/speedtest</b> = Change the current directory to where the main.go file resides.
* <b>go run main.go</b> = Run the program.
* <b>>> cron.txt 2>&1</b> = Dump output from cron job to this file (in the directory we changed to above) and don't mail output.

```0 */4 * * * PATH=$PATH:/usr/local/go/bin && cd /home/pi/Documents/speedtest && go run main.go >> cron.txt 2>&1```
