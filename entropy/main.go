package main

import (
	"fmt"
	"os"

	"github.com/danielchatfield/entropy"
	"github.com/danielchatfield/go-chalk"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(chalk.Red("No file inputted"))
		os.Exit(1)
	}

	filename := os.Args[1]

	rdr, err := entropy.NewFileReaderFromFilename(filename)

	if err != nil {
		fmt.Errorf(chalk.Red("That file doesn't exist"))
		os.Exit(1)
	}

	e := rdr.ShannonEntropy()

	if err != nil {
		fmt.Errorf("An error occured reading the file: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Entropy is  %.2f\n", e)
}
