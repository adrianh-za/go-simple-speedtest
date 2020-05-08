// MIT License

// Copyright (c) 2020 Adrian Houghton

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// speedTestData stores the data of the speedtest
type speedTestData struct {
	URL          string
	Bytes        int
	Milliseconds int
	Date         time.Time
}

// All the download file links are here  -> https://www.thinkbroadband.com/download
// You can choose any of the listed, or any valid file URL you wish.
const fileURL string = "http://ipv4.download.thinkbroadband.com/100MB.zip"

func main() {

	// Do context HTTP download with timeouts
	data, error := speedTestMemoryWithContext(fileURL)
	if error != nil {
		fmt.Println(error)
		return
	}

	// Write data to a file
	error = writeData("speedtest.csv", &data)
	if error != nil {
		fmt.Println(error)
		return
	}
}

// getMegabitsPerSecond returns the Megabits per second (Mbps) of the speedtest data
func getMegabitsPerSecond(data *speedTestData) float64 {
	return getMegabits(data) / getSeconds(data)
}

// getMegabytesPerSecond returns the Megabytes per second (MBps) of the speedtest data
func getMegabytesPerSecond(data *speedTestData) float64 {
	return getMegabytes(data) / getSeconds(data)
}

// getSeconds returns the total seconds of the speedtest data
func getSeconds(data *speedTestData) float64 {
	return float64(data.Milliseconds) / 1000
}

// getMegabytes returns the total Megabytes (MB) of the speedtest data
func getMegabytes(data *speedTestData) float64 {
	return float64(data.Bytes) / 1048576
}

// getMegabits returns the total Megabits (Mb) of the speedtest data
func getMegabits(data *speedTestData) float64 {
	return (float64(data.Bytes) / 1048576) * 8
}

// writeData appends the speedtest data to a file.
func writeData(filename string, data *speedTestData) error {

	// Open the file (Create if not exist)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Build the data string (CSV)
	var dataText strings.Builder
	dataText.WriteString(data.Date.Format("2006-01-02 15:04:05") + ",")
	dataText.WriteString(strconv.Itoa(data.Bytes) + ",")
	dataText.WriteString(strconv.FormatFloat(getMegabytes(data), 'f', 2, 64) + ",")
	dataText.WriteString(strconv.Itoa(data.Milliseconds) + ",")
	dataText.WriteString(strconv.FormatFloat(getSeconds(data), 'f', 2, 64) + ",")
	dataText.WriteString(strconv.FormatFloat(getMegabitsPerSecond(data), 'f', 2, 64) + ",")
	dataText.WriteString(strconv.FormatFloat(getMegabytesPerSecond(data), 'f', 2, 64))

	// Write to the file
	_, err = file.WriteString(dataText.String() + "\n")
	if err != nil {
		return err
	}

	return nil
}

// speedTestMemory downloads the specified file using standard HTTP Client
// A timeout of 120 seconds is set.
func speedTestMemory(url string) (speedTestData, error) {

	// Setup the data struct with defaults
	var data = new(speedTestData)
	data.Milliseconds = -1
	data.Bytes = -1
	data.URL = url
	data.Date = time.Now()

	// Log the START milliseconds
	var startTime = time.Now().UnixNano() / 1000000

	// Setup HTTP client
	var httpClient = &http.Client{
		Timeout: time.Minute * 2,
	}

	// Download the file
	resp, err := httpClient.Get(url)
	if err != nil {
		return *data, err
	}
	defer resp.Body.Close()

	// Read the file into memory
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return *data, err
	}

	// Log the END milliseconds
	var endTime = time.Now().UnixNano() / 1000000

	// Populate the result struct with data and return
	data.Bytes = len(body)
	data.Milliseconds = int(endTime - startTime)
	return *data, nil
}

// speedTestMemory2 downloads the specified file using a context and request
// A timeout of 120 seconds is set.
func speedTestMemoryWithContext(url string) (speedTestData, error) {

	// Setup the data struct with defaults
	var data = new(speedTestData)
	data.Milliseconds = -1
	data.Bytes = -1
	data.URL = url
	data.Date = time.Now()

	// Log the START milliseconds
	var startTime = time.Now().UnixNano() / 1000000

	// Setup HTTP request and context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	// Download the file
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return *data, err
	}
	defer resp.Body.Close()

	// Check if a timeout or cancellation ocurred
	if ctx.Err() != nil {
		return *data, ctx.Err()
	}

	// Read the file into memory
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return *data, err
	}

	// Log the END milliseconds
	var endTime = time.Now().UnixNano() / 1000000

	// Populate the result struct with data and return
	data.Bytes = len(body)
	data.Milliseconds = int(endTime - startTime)
	return *data, nil
}
