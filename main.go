package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
)

const (
	FolderDotConfig       = ".config"
	FolderBookmarkManager = "bookmark-manager"
	FileMainBookmarks     = "main.json"
)

/*
Bookmark Manager: Replaces browser stored bookmark management with a
dedicated solution through an application. This is handy when you migrate
browsers often and don't want to migrate bookmarks every time. All your bookmarks
are in one place.

NOTE: This app assumes, for now, the current user is part of the path to bookmark files.

You can:
  - get a bookmark's URL, name, or description
  - delete a bookmark
  - edit a bookmark's URL, name, or description

from your terminal by using {command} get {bookmark name} and it will
copy the URL to your clipboard.
*/
func main() {
	// Nest structs or just embed bookmark type into a bookmark
	flags := parseFlags()

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error trying to get the current user")
		os.Exit(ExitWithFailure)
	}

	// Make sure the current user's home directory has a .config folder, if not, make it.
	// Repeat that for the bookmark-manager folder
	createDirIfNotExist(path.Join(currentUser.HomeDir, FolderDotConfig))
	createDirIfNotExist(path.Join(currentUser.HomeDir, FolderDotConfig, FolderBookmarkManager))

	mainBookmarksPath := path.Join(currentUser.HomeDir, FolderDotConfig, FolderBookmarkManager, FileMainBookmarks)

	// Create the required main bookmarks file
	createFileIfNotExist(mainBookmarksPath)

	mainFile, fOpenErr := os.OpenFile(mainBookmarksPath, os.O_RDWR, os.ModeAppend)
	defer mainFile.Close()

	if fOpenErr != nil {
		fmt.Println("Error opening main bookmarks file")
		os.Exit(ExitWithFailure)
	}

	bytesRead, fReadErr := os.ReadFile(mainFile.Name())
	if fReadErr != nil {
		fmt.Printf("Unable to read file %s because of error %v\n", mainFile.Name(), fReadErr)
		os.Exit(ExitWithFailure)
	}

	categories := make([]Category, 0)

	// TODO: If bytesRead is empty or nil, you should probably do something else
	unmarshalErr := json.Unmarshal(bytesRead, &categories)
	if unmarshalErr != nil {
		fmt.Printf("Unmarshalling the file's contents caused an error %v\n", unmarshalErr)
	}

	stdin := bufio.NewScanner(os.Stdin)
	handles := Handles{
		Scanner:    stdin,
		File:       mainFile,
		Categories: &categories,
	}

	// TODO: Add a function to list all categories.
	// TODO: Add a function to rename a category.
	if flags.ShouldListAll {
		ListAll(*handles.Categories)
	} else if flags.ShouldEdit {
		EditBookmark(&handles)
	} else if flags.ShouldAdd {
		Add(&handles)
	} else if flags.ShouldDelete {
		Delete(&handles)
	}
}

func parseFlags() CmdFlags {
	shouldListAll := flag.Bool("list", false, "Lists all available bookmarks")
	shouldEdit := flag.Bool("edit", false, "Pass to edit a bookmark. Further prompts will guide you.")
	shouldAdd := flag.Bool("add", false, "Add a new bookmark. Further prompts will guide you.")
	shouldDelete := flag.Bool("delete", false, "Delete an existing bookmark. Further prompts will guide you.")
	flag.Parse()

	return CmdFlags{
		ShouldListAll: *shouldListAll,
		ShouldEdit:    *shouldEdit,
		ShouldAdd:     *shouldAdd,
		ShouldDelete:  *shouldDelete,
	}
}

func createDirIfNotExist(expectedPath string) {
	_, err := os.Stat(expectedPath)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		makeDirErr := os.Mkdir(expectedPath, 0750)
		if makeDirErr != nil {
			fmt.Printf("Unable to make directory: %s\n", expectedPath)
			panic("createDirIfNotExist")
		}
	}
}

func createFileIfNotExist(expectedPath string) {
	fmt.Printf("Requested to create: %s\n", expectedPath)
	_, err := os.Stat(expectedPath)

	if err != nil && errors.Is(err, os.ErrExist) {
		fmt.Println("The file already exists")
		return
	}

	if err != nil && errors.Is(err, os.ErrNotExist) {
		_, createFileErr := os.Create(expectedPath)

		if createFileErr != nil {
			fmt.Printf("Unable to create directory or file %s\n", expectedPath)
			panic(createFileErr)
		}
	}
}
