package main

import (
	"fmt"
	"os"

	"github.com/Pietot/Gonnect-4/grid"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("    --a | --analyse <sequence>")
	fmt.Println("")
	fmt.Println("\033[35m        Analyze a position by giving a score for all possible moves\033[0m")
	fmt.Println("")
	fmt.Println("    --s | --solve   <sequence>")
	fmt.Println("")
	fmt.Println("\033[35m        Solve a position by giving its score and remaining moves\033[0m")
	fmt.Println("")
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
	if len(os.Args) < 3 {
		fmt.Println("\033[31mError: insufficient arguments\033[0m")
		printUsage()
	} else if len(os.Args) > 3 {
		fmt.Println("\033[31mError: too many arguments\033[0m")
		printUsage()
	}

	command := os.Args[1]
	args := os.Args[2:]

	isAnalyze := command == "--a" || command == "--analyse"
	isSolve := command == "--s" || command == "--solve"

	if !isAnalyze && !isSolve {
		fmt.Printf("\033[31mError: unknown command '%s'\033[0m\n", command)
		printUsage()
	}
	if len(os.Args) > 2 {
		for _, a := range os.Args[2:] {
			if a == "--a" || a == "--analyse" || a == "--s" || a == "--solve" {
				fmt.Println("\033[31mError: multiple commands provided\033[0m")
				printUsage()
			}
		}
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

	if isAnalyze {
		fmt.Println(grid.Analyze())
	}

	if isSolve {
		evaluation, stats := grid.Solve()
		fmt.Printf("Evaluation:\n%s\n\nStats: \n%s\n", evaluation, stats)
	}
}
