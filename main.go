package main

import (
	"flag"
)

//
//	Returns None
//
//	Generates a graph of your Git Repositories
func stats(email string) {
	print("stats")
}

func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "Add a new folder to scan for Git Repositories")
	flag.StringVar(&email, "email", "your@email.com", "The email to scan")
	flag.Parse();
	if folder != "" {
		scan(folder)
		return
	}
	stats(email)
}
