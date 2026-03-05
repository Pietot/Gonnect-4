package main

import (
	"fmt"
	"os"

	"github.com/Pietot/Gonnect-4/book"
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
	c.Green("    ./gonnect4.exe -s")
	c.Green("    ./gonnect4.exe -a 32164625")
	c.Green("    ./gonnect4.exe -s --disable-book 5654767662")
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
	// Don't forget to import the book package.
	// Uncomment the following line to continue/recreate building the book.
	book.CreateBook(42)
}
