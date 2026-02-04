package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Pietot/Gonnect-4/config"
	"github.com/Pietot/Gonnect-4/grid"
)

func printUsage() {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("    -a")
	fmt.Println("\033[35m    Analyze a position by giving a score for all possible moves\033[0m")
	fmt.Println()
	fmt.Println("    -s")
	fmt.Println("\033[35m    Solve a position by giving its score and remaining moves\033[0m")
	fmt.Println()
	fmt.Println("    --disable-book")
	fmt.Println("\033[35m    Disable the opening book for analyzing/solving positions\033[0m")
	fmt.Println()
	fmt.Println("    <sequence>")
	fmt.Println("\033[35m    A sequence of digits representing the moves played so far.\n    Columns are numbered from 1 to 7.\033[0m")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("\033[32m    ./gonnect4 -s")
	fmt.Println("    ./gonnect4 -a 32164625")
	fmt.Println("    ./gonnect4 -s --disable-book 5654767662\033[0m")
	fmt.Println()
	os.Exit(1)
}

func parseNumbers(arg string) (string, error) {
	for _, r := range arg {
		if !('0' <= r && r <= '9') {
			return "", fmt.Errorf("\033[31minvalid argument: %s (must contain only digits)\033[0m", arg)
		}
	}
	return arg, nil
}

func main() {
	defer database.DB.Close()
	// Don't forget to import the book package.
	// Uncomment the following line to continue/recreate building the book.
	// book.CreateBook(42)

	flag.Usage = printUsage

	analyze := flag.Bool("a", false, "Analyze a position")
	solve := flag.Bool("s", false, "Solve a position")
	disableBookFlag := flag.Bool("disable-book", false, "Disable opening book")
	config.IsBookEnabled = !*disableBookFlag

	flag.Parse()

	if (*analyze && *solve) || (!*analyze && !*solve) {
		fmt.Println("\033[31mError: you must provide exactly one of -a or -s\033[0m")
		printUsage()
	}

	args := flag.Args()
	if len(args) > 1 {
		fmt.Println("\033[31mError: too many sequences\033[0m")
		printUsage()
	}

	if len(args) == 0 {
		args = append(args, "")
	}

	sequence, err := parseNumbers(args[0])
	if err != nil {
		fmt.Println("\033[31mError:\033[0m", err)
		printUsage()
	}

	grid, err := grid.InitGrid(sequence)
	if err != nil {
		fmt.Println("\033[31mError:\033[0m", err)
		printUsage()
	}

	if *analyze {
		evaluation, stats := grid.Analyze()
		fmt.Printf("Evaluation:\n%s\n\nStats: \n%s\n", evaluation, stats)
	}

	if *solve {
		evaluation, stats := grid.Solve()
		fmt.Printf("Evaluation:\n%s\n\nStats: \n%s\n", evaluation, stats)
	}
}
