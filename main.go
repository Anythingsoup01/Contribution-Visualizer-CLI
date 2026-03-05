package main

import (
	"flag"
	"time"
)



func main() {
	var folder string
	var email string
	var indefinite bool
	var delay_ms int64
	flag.StringVar(&folder, "add", "", "Add a new folder to scan for Git Repositories")
	flag.StringVar(&email, "email", "your@email.com", "The email to scan")
	flag.BoolVar(&indefinite, "i", false, "Loops indefinitely if true")
	flag.Int64Var(&delay_ms, "d", 1000 * 6 * 5, "Set update delay in milliseconds (30000 default)")
	flag.Parse();
	if folder != "" {
		scan(folder)
		return
	}

	delay := time.Duration(delay_ms * int64(time.Millisecond))

	if indefinite {
		for ;;{
			stats(email)
			time.Sleep(delay);
		}
	} else {
		stats(email)
	}
}
