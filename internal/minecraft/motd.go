package minecraft

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/mapstructure"
	"github.com/xrjr/mcutils/pkg/ping"
)

type motd struct {
	Extra []struct {
		Bold   bool   `json:"bold,omitempty"`
		Italic bool   `json:"italic,omitempty"`
		Color  string `json:"color"`
		Text   string `json:"text"`
	} `json:"extra"`
	Text string `json:"text"`
}

var colors = map[string]lipgloss.Color{
	"black":        lipgloss.Color("#000000"),
	"§0":           lipgloss.Color("#000000"),
	"dark_blue":    lipgloss.Color("#0000AA"),
	"§1":           lipgloss.Color("#0000AA"),
	"dark_green":   lipgloss.Color("#00AA00"),
	"§2":           lipgloss.Color("#00AA00"),
	"dark_aqua":    lipgloss.Color("#00AAAA"),
	"§3":           lipgloss.Color("#00AAAA"),
	"dark_red":     lipgloss.Color("#AA0000"),
	"§4":           lipgloss.Color("#AA0000"),
	"dark_purple":  lipgloss.Color("#AA00AA"),
	"§5":           lipgloss.Color("#AA00AA"),
	"gold":         lipgloss.Color("#FFAA00"),
	"§6":           lipgloss.Color("#FFAA00"),
	"gray":         lipgloss.Color("#AAAAAA"),
	"§7":           lipgloss.Color("#AAAAAA"),
	"dark_gray":    lipgloss.Color("#555555"),
	"§8":           lipgloss.Color("#555555"),
	"blue":         lipgloss.Color("#5555FF"),
	"§9":           lipgloss.Color("#5555FF"),
	"green":        lipgloss.Color("#55FF55"),
	"§a":           lipgloss.Color("#55FF55"),
	"aqua":         lipgloss.Color("#55FFFF"),
	"§b":           lipgloss.Color("#55FFFF"),
	"red":          lipgloss.Color("#FF5555"),
	"§c":           lipgloss.Color("#FF5555"),
	"light_purple": lipgloss.Color("#FF55FF"),
	"§d":           lipgloss.Color("#FF55FF"),
	"yellow":       lipgloss.Color("#FFFF55"),
	"§e":           lipgloss.Color("#FFFF55"),
	"white":        lipgloss.Color("#FFFFFF"),
	"§f":           lipgloss.Color("#FFFFFF"),
}

func extra(code string, style *lipgloss.Style) {
	switch code {
	case "§l":
		style.Bold(true)
	case "§o":
		style.Italic(true)
	case "§m":
		style.Strikethrough(true)
	case "§n":
		style.Underline(true)
	case "§r":
		*style = lipgloss.NewStyle()
	}
}

// Motd returns the server Message Of The Day (MOTD) with some formatting.
func Motd(slp ping.JSON) string {
	desc, ok := slp["description"]
	if !ok {
		return ""
	}

	switch typ := desc.(type) {
	case string:
		return motdString(typ)
	case map[string]interface{}:
		return motdMap(typ)
	default:
		return ""
	}
}

func motdString(motd string) string {
	style := lipgloss.NewStyle()
	runes := []rune(motd)
	result := ""

	for i := 0; i < len(runes); i++ {
		char := string(runes[i])
		if char == "§" {
			i++
			mod := string(runes[i])
			char += mod

			if strings.ContainsAny("lomnr", mod) {
				extra(char, &style)
			} else {
				style = style.Foreground(colors[char])
			}
		} else {
			result += style.Render(char)
		}
	}

	return result
}

func motdMap(motdMap map[string]interface{}) string {
	style := lipgloss.NewStyle()
	result := ""

	var msg motd

	mapstructure.Decode(motdMap, &msg)

	for _, ext := range msg.Extra {
		localStyle := style.Copy()
		if ext.Bold {
			localStyle.Bold(true)
		}

		if ext.Italic {
			localStyle.Italic(true)
		}

		if color, ok := colors[ext.Color]; ok {
			localStyle.Foreground(color)
		}

		result += localStyle.Render(ext.Text)
	}

	return result
}
