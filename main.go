package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Pietot/Gonnect-4/config"
	"github.com/Pietot/Gonnect-4/grid"
	c "github.com/fatih/color"
)

func printUsage() {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("    -a")
	c.Magenta("    Analyze a position by giving a score for all possible moves")
	fmt.Println()
	fmt.Println("    -s")
	c.Magenta("    Solve a position by giving its score and remaining moves")
	fmt.Println()
	fmt.Println("    --disable-book")
	c.Magenta("    Disable the opening book for analyzing/solving positions")
	fmt.Println()
	fmt.Println("    <sequence>")
	c.Magenta("    A sequence of digits representing the moves played so far.\n    Columns are numbered from 1 to 7.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	c.Green("    ./gonnect4 -s")
	c.Green("    ./gonnect4 -a 32164625")
	fmt.Println("    ./gonnect4 -s --disable-book 5654767662")
	fmt.Println()
	os.Exit(1)
}

func parseNumbers(arg string) (string, error) {
	for _, r := range arg {
		if !('0' <= r && r <= '9') {
			return "", fmt.Errorf("%s", c.RedString("invalid argument: %s (must contain only digits)", arg))
		}
	}
	return arg, nil
}

func main() {
	flag.Usage = printUsage

	analyze := flag.Bool("a", false, "Analyze a position")
	solve := flag.Bool("s", false, "Solve a position")
	disableBookFlag := flag.Bool("disable-book", false, "Disable opening book")
	flag.Parse()
	
	config.IsBookEnabled = !*disableBookFlag

	if (*analyze && *solve) || (!*analyze && !*solve) {
		c.Red("Error: you must provide exactly one of -a or -s")
		printUsage()
	}

	args := flag.Args()
	if len(args) > 1 {
		c.Red("Error: too many sequences")
		printUsage()
	}

	if len(args) == 0 {
		args = append(args, "")
	}

	sequence, err := parseNumbers(args[0])
	if err != nil {
		c.Red("Error: %s", err)
		printUsage()
	}

	grid, err := grid.InitGrid(sequence)
	if err != nil {
		c.Red("Error: %s", err)
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
