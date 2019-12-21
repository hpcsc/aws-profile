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

	l := widgets.NewList()
	l.Title = "Select a AWS profile"
	l.Rows = labels
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = true
	l.SetRect(0, 0, 100, 15)

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return nil, errors.New("cancelled by user")
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<Enter>":
			return []byte(filteredProfiles[l.SelectedRow]), nil
		}

		ui.Render(l)
	}
}
