package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type column []int

const WEEKS_IN_LAST_SIX_MONTHS = 26;
const DAYS_IN_LAST_SIX_MONTHS = 183;
const OUT_OF_RANGE = 99999;

//
//	Returns time.Time
//
//	Given a time.Time calcualtes the start time of that day
func get_beginning_of_day(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

//
//	Returns int
//
//	Counts how many days since the passed 'date'
func count_days_since_date(date time.Time) int {
	days := 0
	now := get_beginning_of_day(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > DAYS_IN_LAST_SIX_MONTHS {
			return OUT_OF_RANGE
		}
	}
	return days
}

//
//	Returns int
//
//	Calculates the day of week offset,
//	starting at Sunday (7) all the way to Saturday (1)
func calc_offset() int {
	var offset int
	weekday := time.Now().Weekday()
	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}
	return offset
}

//
//	Returns map[int]int
//
//	Fills commits in map, given their path
func fill_commits(email string, path string, commits map[int]int) map[int]int {
	// Instantiate the repo
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		print(path)
		panic(err)
	}
	// Get the commits history starting from HEAD
	itr, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
	// Iterate the commits
	offset := calc_offset()
	err = itr.ForEach(func(c *object.Commit) error {
		daysAgo := count_days_since_date(c.Author.When) + offset
		if (c.Author.Email != email) {
			return nil
		}
		if daysAgo != OUT_OF_RANGE {
			commits[daysAgo]++
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return commits
}

//
//	Returns map[int]int
//
//	Process Repositories from a given user email
func process_repository(email string) map[int]int {
	filePath := get_dot_file_path()
	repos := parse_file_lines_to_slice(filePath)
	daysInMap := DAYS_IN_LAST_SIX_MONTHS
	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}
	for _, path := range repos {
		commits = fill_commits(email, path, commits)
	}

	return commits
}

//
//	Returns []int
//
//	Returns a slice of indexes of the given map
func sort_map_into_slice(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

//
//	Returns map[int]column
//
//	Generates a map with rows and columns ready to print
func build_cols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}
	for _, k := range keys {
		week := int(k/7)
		dayinweek := k % 7
		if dayinweek == 0 {
			col = column{}
		}
		col = append(col, commits[k]);
		if dayinweek == 6 {
			cols[week] = col
		}
	}
	return cols
}

//
//	Returns None
//
//	Prints the month names in the first line
func print_months() {
	week := get_beginning_of_day(time.Now()).Add(-(DAYS_IN_LAST_SIX_MONTHS * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}
		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

//
//	Returns None
//
//	Prints the day of the week
func print_day_col(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}
	fmt.Printf(out)
}

//
//	Returns None
//
//	Prints the value at a given cell as a color
//	with an intesity that goes off the numbers value
func print_cell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}
	if today {
		escape = "\033[1;37;45m"
	}
	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}
	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

//
//	Returns None
//
//	Prints the cells of the graph from the map
func print_cells(cols map[int]column) {
	print_months()
	for j := 6; j >= 0; j-- {
		for i := WEEKS_IN_LAST_SIX_MONTHS + 1; i >= 0; i-- {
			if i == WEEKS_IN_LAST_SIX_MONTHS + 1 {
				print_day_col(j)
			}
			if col, ok := cols[i]; ok {
				if i == 0 && j == calc_offset()-1 {
					print_cell(col[j], true)
					continue
				} else if len(col) > j {
					print_cell(col[j], false)
					continue
				}
			}
			print_cell(0, false)
		}
		fmt.Printf("\n")
	}
}

//
//	Returns None
//
//	Prints the commit stats
func print_commit_stats(commits map[int]int) {
	keys := sort_map_into_slice(commits)
	cols := build_cols(keys, commits)
	print_cells(cols)
}

//
//	Returns None
//
//	Generates a graph of your Git Repositories
func stats(email string) {
	commits := process_repository(email)
	print_commit_stats(commits)
}
