package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkLink(link string) int {
	res, _ := http.Get(link)

	if res != nil {
		res.Body.Close()
		return res.StatusCode
	}

	return 404
}

func main() {
	// Read the entire bookmarks file (exported from Chrome) into memory
	f, err := ioutil.ReadFile("bookmarks.html")
	check(err)

	// Compile the regex for finding links. Very simple because this isn't something
	// that needs to do user validation.
	r, err := regexp.Compile("(?:http|https)://[^'\"]+")
	check(err)

	// Find all of the links in the bookmarks file.
	// The file we read earlier is of []byte type, and the function needs a string so
	// we convert.
	links := r.FindAllString(string(f), -1)

	// Iterate over the returned links.
	for _, l := range links {
		// Make a GET request to the URL and save the status code.
		status := checkLink(l)

		// If we don't get a "200 OK" back and the link isn't for a HTTPS site...
		if status != 200 && !strings.Contains(l, "https") {
			// Check again with a secure connection.
			secureStatus := checkLink(strings.Replace(l, "http", "https", -1))
			// If the link still doesn't return okay...
			if secureStatus != 200 {
				// Then print the link and the status.
				fmt.Println(l + ": " + strconv.Itoa(status))
			}
		}

		// Wait a bit to not spam websites (or make someone's ISP mad).
		time.Sleep(250 * time.Millisecond)
	}
}
