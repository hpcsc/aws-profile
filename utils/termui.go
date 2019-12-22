package utils

import (
	"errors"
	"log"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func filterProfiles(combinedProfiles []string, pattern string) []string {
	var filteredProfiles []string

	for _, profile := range combinedProfiles {
		if pattern == "" || strings.Contains(profile, pattern) {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}

	return filteredProfiles
}

func getDisplayableLabels(profiles []string) []string {
	var labels []string

	for _, profile := range profiles {
		labels = append(labels, strings.Split(profile, ":")[0])
	}

	return labels
}

func SelectProfileFromList(combinedProfiles []string, pattern string) ([]byte, error) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	filteredProfiles := filterProfiles(combinedProfiles, pattern)
	labels := getDisplayableLabels(filteredProfiles)

	list := widgets.NewList()
	list.Title = "Select a AWS profile"
	list.Rows = labels
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
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
			return nil, errors.New("cancelled by user")
		case "j", "<Down>":
			list.ScrollDown()
		case "k", "<Up>":
			list.ScrollUp()
		case "<Enter>":
			return []byte(filteredProfiles[list.SelectedRow]), nil
		}

		ui.Render(grid)
	}
}
