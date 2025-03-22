package main

import (
	"bufio"
	"fmt"
	url2 "net/url"
	"os"
)

const (
	ExitWithFailure = 1
	OSMac           = "darwin"
	OSLinux         = "linux"
	OSWindows       = "windows"
)

type Category struct {
	Type      string     `json:"type"`
	Bookmarks []Bookmark `json:"bookmarks"`
}

func (c Category) String() string {
	return fmt.Sprintf("Category: '%s' has %d bookmarks.", c.Type, len(c.Bookmarks))
}

type Bookmark struct {
	// What do these string literals actually do? can you include the property name in your
	// output?
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	// Tag or folder types should be included
}

// String - Creates a string containing each field from the provided Bookmark. Does not modify any values before
// returning the string.
func (b Bookmark) String() string {
	return fmt.Sprintf("Name %s, URL %s, Description %s", b.Name, b.Url, b.Description)
}

// StringWithTruncatedURL - Creates a string containing each field from the provided Bookmark. Truncates the Bookmark's
// URL to its domain and domain extension (e.g., example.org).
func (b Bookmark) StringWithTruncatedURL() string {
	parsedUrl, urlParseErr := url2.Parse(b.Url)

	printableUrl := ""
	if urlParseErr != nil {
		fmt.Printf("Unable to shorten URL %s, so the entire URL will be used.\n", b.Url)
		printableUrl = b.Url
	}

	// If you care about port number being included in this, swap to .Hostname()
	printableUrl = parsedUrl.Host

	return fmt.Sprintf("Name %s, URL %s, Description %s", b.Name, printableUrl, b.Description)
}

type CmdFlags struct {
	ShouldListAll bool
	ShouldEdit    bool
	ShouldAdd     bool
	ShouldDelete  bool
	ShouldSelect  bool
}

type Handles struct {
	Scanner    *bufio.Scanner
	Categories *[]Category
	File       *os.File
}

type Indexes struct {
	CategoryIndex int
	BookmarkIndex int
}

var BookmarkNotFoundError = fmt.Errorf("bookmark not found in categories")
