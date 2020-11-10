package tui

import (
	"github.com/hpcsc/aws-profile/internal/utils"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/hpcsc/aws-profile/internal/config"
)

func getDisplayableLabels(profiles []config.Profile) []string {
	var labels []string

	for _, profile := range profiles {
		labels = append(labels, profile.DisplayProfileName)
	}

	return labels
}

func SelectProfileFromList(profiles config.Profiles, pattern string) ([]byte, error) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
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
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	selectedIndex, err := renderListSelection(values, title)
	if err != nil {
		return nil, err
	}
	return []byte(values[selectedIndex]), nil
}

func renderListSelection(labels []string, title string) (int, error) {
	list := widgets.NewList()
	list.Title = title
	list.Rows = labels
	list.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
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
