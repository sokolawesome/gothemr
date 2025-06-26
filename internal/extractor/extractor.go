package extractor

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Color struct {
	R, G, B uint8
	Count   int
}

func (c Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func (c Color) RGB() (uint8, uint8, uint8) {
	return c.R, c.G, c.B
}

func (c Color) TerminalString() string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm█\x1b[0m", c.R, c.G, c.B)
}

func ExtractColors(imagePath string, count int) ([]Color, error) {
	img, err := loadImage(imagePath)
	if err != nil {
		return nil, err
	}

	return extractKMeans(img, count)
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return png.Decode(file)
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	case ".gif":
		gifImg, err := gif.DecodeAll(file)
		if err != nil {
			return nil, err
		}
		if len(gifImg.Image) == 0 {
			return nil, fmt.Errorf("empty gif")
		}
		middleFrame := len(gifImg.Image) / 2
		return gifImg.Image[middleFrame], nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
}

func extractKMeans(img image.Image, count int) ([]Color, error) {
	bounds := img.Bounds()
	var pixels []Color

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 3 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 3 {
			r, g, b, a := img.At(x, y).RGBA()
			if a < 128 {
				continue
			}

			color := quantizeColor(Color{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
			}, 16)

			if isValidColor(color) {
				pixels = append(pixels, color)
			}
		}
	}

	if len(pixels) == 0 {
		return nil, fmt.Errorf("no valid pixels found")
	}

	centroids := kmeansCluster(pixels, count)

	sort.Slice(centroids, func(i, j int) bool {
		return centroids[i].Count > centroids[j].Count
	})

	return centroids, nil
}

func kmeansCluster(pixels []Color, k int) []Color {
	if len(pixels) <= k {
		return pixels
	}

	centroids := make([]Color, k)
	for i := range k {
		centroids[i] = pixels[i*len(pixels)/k]
	}

	for range 10 {
		clusters := make([][]Color, k)

		for _, p := range pixels {
			closest := 0
			minDist := colorDistance(p, centroids[0])
			for i := range k {
				if dist := colorDistance(p, centroids[i]); dist < minDist {
					minDist = dist
					closest = i
				}
			}
			clusters[closest] = append(clusters[closest], p)
		}

		changed := false
		for i := range k {
			if len(clusters[i]) == 0 {
				centroids[i] = pixels[rand.Intn(len(pixels))]
				changed = true
				continue
			}
			newCentroid := averageColor(clusters[i])
			if newCentroid != centroids[i] {
				centroids[i] = newCentroid
				changed = true
			}
		}

		if !changed {
			break
		}
	}

	counts := make([]int, k)
	for _, p := range pixels {
		best := 0
		minDist := colorDistance(p, centroids[0])
		for i := range k {
			if dist := colorDistance(p, centroids[i]); dist < minDist {
				best = i
				minDist = dist
			}
		}
		counts[best]++
	}

	result := make([]Color, k)
	for i, c := range centroids {
		result[i] = Color{
			R:     c.R,
			G:     c.G,
			B:     c.B,
			Count: counts[i],
		}
	}
	return result
}

func colorDistance(c1, c2 Color) float64 {
	dr := float64(c1.R) - float64(c2.R)
	dg := float64(c1.G) - float64(c2.G)
	db := float64(c1.B) - float64(c2.B)
	return 0.3*dr*dr + 0.59*dg*dg + 0.11*db*db
}

func averageColor(colors []Color) Color {
	var r, g, b int
	for _, c := range colors {
		r += int(c.R)
		g += int(c.G)
		b += int(c.B)
	}
	n := len(colors)
	return Color{
		R: uint8(r / n),
		G: uint8(g / n),
		B: uint8(b / n),
	}
}

func quantizeColor(c Color, step uint8) Color {
	return Color{
		R: (c.R / step) * step,
		G: (c.G / step) * step,
		B: (c.B / step) * step,
	}
}

func isValidColor(c Color) bool {
	brightness := (int(c.R) + int(c.G) + int(c.B)) / 3
	return brightness > 20 && brightness < 235
}
