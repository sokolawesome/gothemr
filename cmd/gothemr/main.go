package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sokolawesome/gothemr/internal/config"
	"github.com/sokolawesome/gothemr/internal/extractor"
	"github.com/sokolawesome/gothemr/internal/palette"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Config loaded. Cache directory: %s\n\n", cfg.CacheDir)

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

	pal := palette.Generate(colors)

	fmt.Println("--- Palette ---")
	fmt.Printf("Background: %s\n", pal.Background.Hex())
	fmt.Printf("Foreground: %s\n", pal.Foreground.Hex())
	fmt.Printf("Accent:     %s\n", pal.Accent.Hex())
}
