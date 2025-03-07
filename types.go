package main

import (
	"bufio"
	"fmt"
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

func (b Bookmark) String() string {
	return fmt.Sprintf("Name: %s, Url: %s, Description: %s", b.Name, b.Url, b.Description)
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
