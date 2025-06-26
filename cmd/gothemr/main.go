package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sokolawesome/gothemr/internal/extractor"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide an image path.")
	}
	imagePath := os.Args[1]

	colors, err := extractor.ExtractColors(imagePath, 16)
	if err != nil {
		log.Fatalf("Failed to extract colors: %v", err)
	}

	fmt.Println("Extracted Colors:")
	for _, color := range colors {
		fmt.Printf("%s %s (R:%d G:%d B:%d)\n", color.TerminalString(), color.Hex(), color.R, color.G, color.B)
	}
}
