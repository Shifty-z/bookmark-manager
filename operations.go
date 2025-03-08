package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	URLRequiredPrefixProtocol = "https://"
	URLRequiredPrefixWWWDot   = "www."
)

func ListAll(categories []Category) {
	fmt.Print("Listing all available categories and bookmarks\n\n")

	for catIdx, value := range categories {
		// Skip empty categories
		if len(categories[catIdx].Bookmarks) == 0 {
			continue
		}

		fmt.Printf("%d %s\n", catIdx, value.String())

		for bookmarkIdx, bookmark := range value.Bookmarks {
			fmt.Printf("--> %d %s\n", bookmarkIdx, bookmark.String())
		}
		fmt.Println("")
	}
}

func EditBookmark(handles *Handles) {
	fmt.Print("Editing available bookmarks. What is the bookmark's name? ")
	providedName := getScannedInput(handles.Scanner)

	indexes, notFoundErr := findBookmarkByName(handles, providedName)
	if notFoundErr != nil {
		fmt.Printf("No bookmark with the name '%s' was found.\n", providedName)
		os.Exit(ExitWithFailure)
	}

	fmt.Println("If you want to leave a value set to what it was, press enter. Otherwise, enter a new value")
	fmt.Printf("Update Bookmark named '%s': ", (*handles.Categories)[indexes.CategoryIndex].Bookmarks[indexes.BookmarkIndex].Name)
	updatedName := getScannedInput(handles.Scanner)
	if updatedName != "" {
		(*handles.Categories)[indexes.CategoryIndex].Bookmarks[indexes.BookmarkIndex].Name = updatedName
	}

	fmt.Printf("Update the URL '%s': ", (*handles.Categories)[indexes.CategoryIndex].Bookmarks[indexes.BookmarkIndex].Url)
	updatedUrl := getScannedInput(handles.Scanner)
	if updatedUrl != "" {
		(*handles.Categories)[indexes.CategoryIndex].Bookmarks[indexes.BookmarkIndex].Url = updatedUrl
	}

	fmt.Printf("Update the description '%s': ", (*handles.Categories)[indexes.CategoryIndex].Bookmarks[indexes.BookmarkIndex].Description)
	updatedDescription := getScannedInput(handles.Scanner)
	if updatedDescription != "" {
		(*handles.Categories)[indexes.CategoryIndex].Bookmarks[indexes.BookmarkIndex].Description = updatedDescription
	}

	writeFile(handles)
}

func AddBookmark(handles *Handles) {
	// TODO: Handle uncategorized bookmarks if you don't want to supply a category
	fmt.Print("Enter a category type for this bookmark: ")
	category := getScannedInput(handles.Scanner)

	fmt.Print("Enter a name for the bookmark: ")
	newName := getScannedInput(handles.Scanner)

	fmt.Print("Enter a URL for the bookmark: ")
	newUrl := getScannedInput(handles.Scanner)

	protocolHttp := "http://"
	hasHttpProtocol := strings.HasPrefix(newUrl, protocolHttp)
	if hasHttpProtocol {
		protocolLength := len(protocolHttp)
		newUrl = URLRequiredPrefixProtocol + newUrl[protocolLength:]
	}

	hasHttpsProtocol := strings.HasPrefix(newUrl, URLRequiredPrefixProtocol)
	hasWwwDot := strings.Contains(newUrl, URLRequiredPrefixWWWDot)
	if !hasHttpsProtocol && !hasWwwDot {
		newUrl = URLRequiredPrefixProtocol + URLRequiredPrefixWWWDot + newUrl
	} else if !hasHttpsProtocol && hasWwwDot {
		newUrl = URLRequiredPrefixProtocol + newUrl
	} else if hasHttpsProtocol && !hasWwwDot {
		httpsProtocolColonTwoSlashesLength := 8
		newUrl = URLRequiredPrefixProtocol + URLRequiredPrefixWWWDot + newUrl[httpsProtocolColonTwoSlashesLength:]
	}

	fmt.Print("Describe the bookmark. What's it for? ")
	newDesc := getScannedInput(handles.Scanner)

	fmt.Printf("Category: %s\n", category)
	fmt.Printf("Name: %s\n", newName)
	fmt.Printf("Url: %s\n", newUrl)
	fmt.Printf("Desc: %s\n", newDesc)

	newBookmark := Bookmark{
		Name:        newName,
		Url:         newUrl,
		Description: newDesc,
	}

	// If the category provided exists, append to that.
	wasCatFound := false
	for catIdx, cat := range *handles.Categories {
		if category == cat.Type {
			fmt.Println("Category was found.")
			(*handles.Categories)[catIdx].Bookmarks = append((*handles.Categories)[catIdx].Bookmarks, newBookmark)
			wasCatFound = true
		}
	}

	// Create a new category because this one has never been supplied
	if !wasCatFound {
		fmt.Println("Category was not found. Creating a new one.")
		newCategory := Category{
			Type:      category,
			Bookmarks: []Bookmark{newBookmark},
		}
		*handles.Categories = append(*handles.Categories, newCategory)
	}

	writeFile(handles)
}

func DeleteBookmark(handles *Handles) {
	fmt.Print("Provide a name for the bookmark you want to delete: ")
	name := getScannedInput(handles.Scanner)
	fmt.Printf("Searching for bookmark '%s'\n", name)

	indexes, notFoundErr := findBookmarkByName(handles, name)
	if notFoundErr != nil {
		fmt.Printf("No bookmark with the name '%s' was found.\n", name)
		os.Exit(ExitWithFailure)
	}

	// Make a slice from 0 to index, then index + 1 to end of arr
	updatedBookmarks := append((*handles.Categories)[indexes.CategoryIndex].Bookmarks[:indexes.BookmarkIndex], (*handles.Categories)[indexes.CategoryIndex].Bookmarks[indexes.BookmarkIndex+1:]...)

	(*handles.Categories)[indexes.CategoryIndex].Bookmarks = updatedBookmarks
	fmt.Printf("Bookmark: %s has been removed from category %s\n", name, (*handles.Categories)[indexes.CategoryIndex].Type)

	writeFile(handles)
}

func SelectBookmark(handles *Handles) {
	fmt.Print("Enter the category number the bookmark is in: ")
	categoryInput := getScannedInput(handles.Scanner)
	fmt.Printf("You entered category '%s'\n", categoryInput)

	fmt.Print("Enter the bookmark number: ")
	bookmarkInput := getScannedInput(handles.Scanner)
	fmt.Printf("You entered bookmark number: '%s'\n", bookmarkInput)

	categoryNumber, catConvertErr := strconv.Atoi(categoryInput)
	bookmarkNumber, bookmarkConvertErr := strconv.Atoi(bookmarkInput)

	if catConvertErr != nil {
		fmt.Printf("Could not convert input while selecting a category. Error: %v\n", catConvertErr)
		os.Exit(ExitWithFailure)
	}

	if bookmarkConvertErr != nil {
		fmt.Printf("Could not convert input while selecting a bookmark. Error: %v\n", bookmarkConvertErr)
		os.Exit(ExitWithFailure)
	}

	fmt.Printf("Selected Category '%s' bookmark named %s\n", (*handles.Categories)[categoryNumber].Type, (*handles.Categories)[categoryNumber].Bookmarks[bookmarkNumber].Name)
	selectedURL := (*handles.Categories)[categoryNumber].Bookmarks[bookmarkNumber].Url

	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case OSMac:
		startErr := exec.Command("open", selectedURL).Start()
		if startErr != nil {
			fmt.Printf("Error opening the browser to %s because of error: %v\n", selectedURL, startErr)
			os.Exit(ExitWithFailure)
		}
	case OSLinux:
		fmt.Println("Linux is not supported, yet.")
		os.Exit(ExitWithFailure)
	case OSWindows:
		fmt.Println("Windows is not supported, yet.")
		os.Exit(ExitWithFailure)
	}
}

func writeFile(handles *Handles) {
	clearFileDataErr := handles.File.Truncate(0)
	// Not sure if this is needed.
	//handles.File.Seek(0, 0)

	if clearFileDataErr != nil {
		fmt.Println("Error trying to truncate file while editing its contents.")
		panic(clearFileDataErr)
	}

	writableBookmarks := marshalBytes(handles.Categories)
	_, writeErr := handles.File.Write(writableBookmarks)

	if writeErr != nil {
		fmt.Printf("Unable to write all bytes to file %s! To prevent data loss, I will dump contents to standard in!\n", handles.File.Name())
		// TODO: Test this
		fmt.Printf("%v\n", *handles.Categories)
	}
}

func marshalBytes(bookmarks *[]Category) []byte {
	convertedBytes, marshalErr := json.Marshal(*bookmarks)
	if marshalErr != nil {
		panic(fmt.Sprintf("Problem marhshalling bytes: %s\n", marshalErr))
	}

	return convertedBytes
}

func getScannedInput(scanner *bufio.Scanner) string {
	scanner.Scan()
	input := scanner.Text()
	return input
}

func doesBookmarkExist(categories []Category, bookmarkName string) bool {
	for _, cat := range categories {
		for _, bookmark := range cat.Bookmarks {
			if bookmark.Name == bookmarkName {
				return true
			}
		}
	}

	return false
}

func findBookmarkByName(handles *Handles, bookmarkName string) (Indexes, error) {
	doesExist := doesBookmarkExist(*handles.Categories, bookmarkName)
	if !doesExist {
		return Indexes{}, fmt.Errorf("Error: %w! Searched for name: %s\n", BookmarkNotFoundError, bookmarkName)
	}

	// Collect all bookmarks the same names; Allow yourself to select the exact bookmark you want to edit
	bookmarkIndexes := make([]Indexes, 0)
	for catIdx, cat := range *handles.Categories {
		for bookmarkIdx, bookmark := range cat.Bookmarks {
			if bookmark.Name == bookmarkName {
				bookmarkIndexes = append(bookmarkIndexes, Indexes{
					CategoryIndex: catIdx,
					BookmarkIndex: bookmarkIdx,
				})
			}
		}
	}

	hasMultipleMatches := len(bookmarkIndexes) > 1
	if hasMultipleMatches {
		fmt.Printf("Multiple bookmarks with the name %s were found.\n", bookmarkName)
		for idx, bookmark := range bookmarkIndexes {
			fmt.Printf("Category '%s'\n", (*handles.Categories)[bookmark.CategoryIndex].Type)
			fmt.Printf("--> %d %s\n", idx, (*handles.Categories)[bookmark.CategoryIndex].Bookmarks[bookmark.BookmarkIndex].String())
		}

		fmt.Printf("Which number corresponds to the bookmark you'd like to select? ")
		input := getScannedInput(handles.Scanner)
		selectedIndex, convertErr := strconv.Atoi(input)

		if convertErr != nil {
			fmt.Printf("Unable to conver your input into an integer type. %s", convertErr)
			os.Exit(ExitWithFailure)
		}

		fmt.Printf("Entered {%d} as the selected index. BM Indexes length: {%d}\n", selectedIndex, len(bookmarkIndexes))
		if selectedIndex < 0 || selectedIndex >= len(bookmarkIndexes) {
			fmt.Println("Your input is out of the accepted boundaries. Try a number that was listed.")
			os.Exit(ExitWithFailure)
		}

		return bookmarkIndexes[selectedIndex], nil
	}

	return bookmarkIndexes[0], nil
}
