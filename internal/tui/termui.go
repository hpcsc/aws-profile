package tui

import (
	"fmt"
	"github.com/hpcsc/aws-profile/internal/config"
	"github.com/hpcsc/aws-profile/internal/utils"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/hpcsc/aws-profile/internal/awsconfig"
)

func getDisplayableLabels(profiles []awsconfig.Profile) []string {
	var labels []string

	for _, profile := range profiles {
		labels = append(labels, profile.DisplayProfileName)
	}

	return labels
}

func SelectProfileFromList(profiles awsconfig.Profiles, pattern string) ([]byte, error) {
	if err := ui.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	filteredProfiles := profiles.Filter(pattern)
	labels := getDisplayableLabels(filteredProfiles)

	selectedIndex, err := renderListSelection(labels, "Select an AWS profile")
	if err != nil {
		return nil, err
	}
	return []byte(filteredProfiles[selectedIndex].ProfileName), nil
}

func SelectValueFromList(values []string, title string) ([]byte, error) {
	if err := ui.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	selectedIndex, err := renderListSelection(values, title)
	if err != nil {
		return nil, err
	}
	return []byte(values[selectedIndex]), nil
}

func toTermUIColor(color string) ui.Color {
	switch strings.ToLower(color) {
	case "black":
		return ui.ColorBlack
	case "red":
		return ui.ColorRed
	case "yellow":
		return ui.ColorYellow
	case "blue":
		return ui.ColorBlue
	case "magenta":
		return ui.ColorMagenta
	case "cyan":
		return ui.ColorCyan
	case "white":
		return ui.ColorWhite
	}

	return ui.ColorGreen
}

func renderListSelection(labels []string, title string) (int, error) {
	c, err := config.Load()
	if err != nil {
		return -1, err
	}

	list := widgets.NewList()
	list.Title = title
	list.Rows = labels
	list.SelectedRowStyle = ui.NewStyle(toTermUIColor(c.HighlightColor))
	list.WrapText = true

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/3,
			ui.NewCol(1.0, list),
		),
	)

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return -1, utils.NewCancelledError()
		case "j", "<Down>":
			list.ScrollDown()
		case "k", "<Up>":
			list.ScrollUp()
		case "<Enter>":
			return list.SelectedRow, nil
		}

		ui.Render(grid)
	}
}
