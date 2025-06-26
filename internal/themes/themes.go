package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sokolawesome/gothemr/internal/extractor"
	"github.com/sokolawesome/gothemr/internal/palette"
)

func GenerateAll(pal *palette.Palette, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	generators := map[string]func(*palette.Palette) (string, error){
		"hyprland.conf":     generateHyprland,
		"waybar.css":        generateWaybar,
		"rofi.rasi":         generateRofi,
		"gtk.css":           generateGTK,
		"colors-kitty.conf": generateKitty,
	}

	for filename, generator := range generators {
		content, err := generator(pal)
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}

		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

func generateHyprland(pal *palette.Palette) (string, error) {
	var sb strings.Builder
	sb.WriteString("general {\n")
	sb.WriteString(fmt.Sprintf("    col.active_border = rgb(%s)\n", strings.TrimPrefix(pal.Accent.Hex(), "#")))
	sb.WriteString(fmt.Sprintf("    col.inactive_border = rgb(%s)\n", strings.TrimPrefix(pal.Background.Hex(), "#")))
	sb.WriteString("}\n\n")

	sb.WriteString("decoration {\n")
	sb.WriteString(fmt.Sprintf("    col.shadow = rgb(%s)\n", strings.TrimPrefix(pal.Background.Hex(), "#")))
	sb.WriteString("}\n\n")

	sb.WriteString("group {\n")
	sb.WriteString(fmt.Sprintf("    col.border_active = rgb(%s)\n", strings.TrimPrefix(pal.Accent.Hex(), "#")))
	sb.WriteString(fmt.Sprintf("    col.border_inactive = rgb(%s)\n", strings.TrimPrefix(pal.Background.Hex(), "#")))
	sb.WriteString("}\n")

	return sb.String(), nil
}

func generateWaybar(pal *palette.Palette) (string, error) {
	var sb strings.Builder
	sb.WriteString("@define-color foreground " + pal.Foreground.Hex() + ";\n")
	sb.WriteString("@define-color background " + pal.Background.Hex() + ";\n")
	sb.WriteString("@define-color cursor " + pal.Accent.Hex() + ";\n")

	for i, color := range pal.Colors {
		sb.WriteString(fmt.Sprintf("@define-color color%d %s;\n", i, color.Hex()))
	}

	return sb.String(), nil
}

func generateRofi(pal *palette.Palette) (string, error) {
	var sb strings.Builder
	sb.WriteString("* {\n")
	sb.WriteString(fmt.Sprintf("    background: %s;\n", pal.Background.Hex()))
	sb.WriteString(fmt.Sprintf("    foreground: %s;\n", pal.Foreground.Hex()))
	sb.WriteString(fmt.Sprintf("    selected: %s;\n", pal.Accent.Hex()))
	sb.WriteString(fmt.Sprintf("    active: %s;\n", pal.Cursor.Hex()))
	sb.WriteString("}\n\n")

	sb.WriteString("window {\n")
	sb.WriteString("    transparency: \"real\";\n")
	sb.WriteString("    background-color: @background;\n")
	sb.WriteString("    border: 2px;\n")
	sb.WriteString("    border-color: @selected;\n")
	sb.WriteString("}\n\n")

	sb.WriteString("mainbox {\n")
	sb.WriteString("    background-color: @background;\n")
	sb.WriteString("}\n\n")

	sb.WriteString("element selected {\n")
	sb.WriteString("    background-color: @selected;\n")
	sb.WriteString("    text-color: @background;\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}

func generateGTK(pal *palette.Palette) (string, error) {
	var sb strings.Builder
	sb.WriteString("@define-color theme_bg_color " + pal.Background.Hex() + ";\n")
	sb.WriteString("@define-color theme_fg_color " + pal.Foreground.Hex() + ";\n")
	sb.WriteString("@define-color theme_selected_bg_color " + pal.Accent.Hex() + ";\n")
	sb.WriteString("@define-color theme_selected_fg_color " + pal.Background.Hex() + ";\n")
	sb.WriteString("@define-color insensitive_bg_color " + pal.Special["medium_dark"].Hex() + ";\n")
	sb.WriteString("@define-color insensitive_fg_color " + pal.Special["medium"].Hex() + ";\n")
	sb.WriteString("@define-color borders " + pal.Special["dark"].Hex() + ";\n")
	sb.WriteString("@define-color warning_color " + pal.Special["yellow"].Hex() + ";\n")
	sb.WriteString("@define-color error_color " + pal.Special["red"].Hex() + ";\n")
	sb.WriteString("@define-color success_color " + pal.Special["green"].Hex() + ";\n")

	return sb.String(), nil
}

func generateKitty(pal *palette.Palette) (string, error) {
	var sb strings.Builder

	return sb.String(), nil
}

func adjustBrightness(color extractor.Color, factor float64) extractor.Color {
	r := float64(color.R) * factor
	g := float64(color.G) * factor
	b := float64(color.B) * factor

	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}

	return extractor.Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}
}
