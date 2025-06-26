package palette

import (
	"sort"

	"github.com/sokolawesome/gothemr/internal/extractor"
)

type Palette struct {
	Colors     []extractor.Color
	Background extractor.Color
	Foreground extractor.Color
	Cursor     extractor.Color
	Accent     extractor.Color
	Special    map[string]extractor.Color
}

func Generate(colors []extractor.Color) *Palette {
	if len(colors) == 0 {
		return &Palette{}
	}

	sortedColors := make([]extractor.Color, len(colors))
	copy(sortedColors, colors)
	sort.Slice(sortedColors, func(i, j int) bool {
		return getBrightness(sortedColors[i]) < getBrightness(sortedColors[j])
	})

	pal := &Palette{
		Colors:  colors,
		Special: make(map[string]extractor.Color),
	}

	pal.Background = sortedColors[0]
	pal.Foreground = sortedColors[len(sortedColors)-1]

	if len(sortedColors) > 2 {
		pal.Cursor = sortedColors[len(sortedColors)-2]
	} else {
		pal.Cursor = pal.Foreground
	}

	pal.Accent = findAccentColor(colors)

	pal.generateSpecialColors()

	return pal
}

func getBrightness(c extractor.Color) int {
	return int(c.R)*299 + int(c.G)*587 + int(c.B)*114
}

func findAccentColor(colors []extractor.Color) extractor.Color {
	if len(colors) == 0 {
		return extractor.Color{}
	}

	maxSaturation := 0.0
	var accentColor extractor.Color

	for _, color := range colors {
		saturation := getSaturation(color)
		if saturation > maxSaturation {
			maxSaturation = saturation
			accentColor = color
		}
	}

	return accentColor
}

func getSaturation(c extractor.Color) float64 {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := r
	if g > max {
		max = g
	}
	if b > max {
		max = b
	}

	min := r
	if g < min {
		min = g
	}
	if b < min {
		min = b
	}

	if max == 0 {
		return 0
	}

	return (max - min) / max
}

func (p *Palette) generateSpecialColors() {
	if len(p.Colors) < 8 {
		return
	}

	sortedByBrightness := make([]extractor.Color, len(p.Colors))
	copy(sortedByBrightness, p.Colors)
	sort.Slice(sortedByBrightness, func(i, j int) bool {
		return getBrightness(sortedByBrightness[i]) < getBrightness(sortedByBrightness[j])
	})

	third := len(sortedByBrightness) / 3

	p.Special["dark"] = sortedByBrightness[0]
	p.Special["medium_dark"] = sortedByBrightness[third]
	p.Special["medium"] = sortedByBrightness[third*2]
	p.Special["light"] = sortedByBrightness[len(sortedByBrightness)-1]

	redColors := filterByHue(p.Colors, 0, 30)
	if len(redColors) > 0 {
		p.Special["red"] = redColors[0]
	}

	greenColors := filterByHue(p.Colors, 90, 150)
	if len(greenColors) > 0 {
		p.Special["green"] = greenColors[0]
	}

	blueColors := filterByHue(p.Colors, 200, 260)
	if len(blueColors) > 0 {
		p.Special["blue"] = blueColors[0]
	}

	yellowColors := filterByHue(p.Colors, 45, 75)
	if len(yellowColors) > 0 {
		p.Special["yellow"] = yellowColors[0]
	}
}

func filterByHue(colors []extractor.Color, minHue, maxHue float64) []extractor.Color {
	var filtered []extractor.Color

	for _, color := range colors {
		hue := getHue(color)
		if hue >= minHue && hue <= maxHue {
			filtered = append(filtered, color)
		}
	}

	return filtered
}

func getHue(c extractor.Color) float64 {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	max := r
	if g > max {
		max = g
	}
	if b > max {
		max = b
	}

	min := r
	if g < min {
		min = g
	}
	if b < min {
		min = b
	}

	if max == min {
		return 0
	}

	var hue float64
	switch max {
	case r:
		hue = (g - b) / (max - min)
	case g:
		hue = 2.0 + (b-r)/(max-min)
	case b:
		hue = 4.0 + (r-g)/(max-min)
	}

	hue *= 60
	if hue < 0 {
		hue += 360
	}

	return hue
}
