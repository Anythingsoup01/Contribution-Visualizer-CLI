package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"slices"
	"strings"
)

//
//	Returns []string of relative .git folder paths
//
//	Recursively searches in the subfolders by passing an existing 'Folders' slice.
func scan_git_folders(folders []string, folder string) []string {
	// Trims the last '/'
	folder = strings.TrimSuffix(folder, "/")

	
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	var path string
	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Print(path + "\n")
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" {
				continue
			}
			folders = scan_git_folders(folders, path)
		}
	}
	return folders
}

//
//	Returns []string of folders
//
//	Helper functions to call scan_git_folders cleanly
func recursive_scan_folder(folder string) []string {
	return scan_git_folders(make([]string, 0), folder)
}

//
//	Returns string path to .gogitlocalstats
//
func get_dot_file_path() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dotFile := usr.HomeDir + "/.gogitlocalstats"
	return dotFile
}

//
//	Returns *os.File
//
//	Opens a file and creates it if it doesn't exist
func open_file(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND | os.O_RDWR, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return f
}

//
//	Returns []string
//
//	Parses a file and returns it's contents as blocks of lines
func parse_file_lines_to_slice(filePath string) []string {
	f := open_file(filePath)
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			print("Error")
			panic(err)
		}
	}

	return lines

}

//
//	Returns boolean
//
//	Checks if a slice exists in an already existing slice stack.
//	Returns true if the slice is found, false otherwise
func slice_contains(slice []string , value string) bool {
	return slices.Contains(slice, value)
}

//
//	Returns []string
//
//	Joins new slices to already existing slices
func join_slices(new []string, existing []string) []string {
	for _, i := range new {
		if !slice_contains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func dump_strings_slice_to_file(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(content), 0755)
}

//
//	Returns None
//
//	Stores the given slice to the filesystem
func add_new_slice_elements_to_file(filePath string, newRepos []string) {
	existingRepos := parse_file_lines_to_slice(filePath)
	repos := join_slices(newRepos, existingRepos)
	dump_strings_slice_to_file(repos, filePath)
}

//
//	Returns None
//
//	Scans a given path and crawls through subfolders
//	to find Git Repositories
func scan(folder string) {
	fmt.Printf("Found Folders:\n\n")
	repositories := recursive_scan_folder(folder)
	filePath := get_dot_file_path()
	add_new_slice_elements_to_file(filePath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n");
}
