package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func logo() string {
	colors := [][]int{
		{114, 180, 74},
		{99, 163, 67},
		{116, 184, 75},
		{104, 166, 67},
		{83, 144, 64},
		{146, 196, 99},

		{146, 196, 99},
		{88, 59, 41},
		{104, 166, 67},
		{83, 144, 64},
		{116, 184, 75},
		{121, 87, 59},

		{88, 59, 41},
		{121, 87, 59},
		{184, 135, 92},
		{151, 108, 74},
		{104, 166, 67},
		{88, 59, 41},

		{184, 135, 92},
		{136, 138, 137},
		{151, 108, 74},
		{184, 135, 92},
		{88, 59, 41},
		{121, 87, 59},

		{151, 108, 74},
		{184, 135, 92},
		{184, 135, 92},
		{88, 59, 41},
		{121, 87, 59},
		{184, 135, 92},

		{184, 135, 92},
		{88, 59, 41},
		{184, 135, 92},
		{151, 108, 74},
		{136, 138, 137},
		{121, 87, 59},
	}
	logo := ""

	for i, color := range colors {
		logo += fmt.Sprintf("\033[38;2;%d;%d;%dm", color[0], color[1], color[2])
		logo += "██"
		logo += "\033[0m"

		if (i+1)%6 == 0 {
			logo += "\n"
		}
	}

	return logo
}

func header() string {
	return `███╗   ███╗██╗███╗   ██╗███████╗███████╗███████╗███████╗██╗  ██╗
████╗ ████║██║████╗  ██║██╔════╝██╔════╝██╔════╝██╔════╝██║ ██╔╝
██╔████╔██║██║██╔██╗ ██║█████╗  ███████╗█████╗  █████╗  █████╔╝
██║╚██╔╝██║██║██║╚██╗██║██╔══╝  ╚════██║██╔══╝  ██╔══╝  ██╔═██╗
██║ ╚═╝ ██║██║██║ ╚████║███████╗███████║███████╗███████╗██║  ██╗
╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝`
}

func printLogo() {
	w := lipgloss.Width(header())

	fmt.Println()
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().MarginRight(5).Render(logo()),
		lipgloss.JoinVertical(lipgloss.Top,
			header(),
			lipgloss.PlaceHorizontal(w, lipgloss.Center, fmt.Sprintf("%s - MrMarble", version)),
		),
	))
	fmt.Println()
}
