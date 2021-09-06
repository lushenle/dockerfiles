package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
)

// difference returns the elements in `slice1` that aren't in `slice2`
func difference(slice1 []string, slice2 []string) []string {
	var result []string

	// Loop two times, first to find s strings not in s2,
	// second loop to find s2 strings not in s
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found, add it to return slice
			if !found {
				result = append(result, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return result
}

// readLines read the records in the file into a slice
func readLines(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result []string

	// Scan the entire file
	scanner := bufio.NewScanner(f)
	// Read line by line, add line to slice
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	lines, _ := readLines(os.Args[1])
	lines2, _ := readLines(os.Args[2])

	wg.Add(1)
	go func() {
		for _, v := range difference(lines, lines2) {
			fmt.Println(v)
		}
		defer wg.Done()
	}()
	wg.Wait()
}

/*
func main() {
	lines := make(map[string]bool)

	f, err := os.Open(os.Args[1])
	fsc := bufio.NewScanner(f)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	for fsc.Scan() {
		lines[fsc.Text()] = true
	}

	f2, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	defer f2.Close()

	f2sc := bufio.NewScanner(f2)

	for f2sc.Scan() {
		if !lines[f2sc.Text()] {
			fmt.Println(f2sc.Text())
		}
	}
}
*/
