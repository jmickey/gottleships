package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

// Load ...loads the UI...
func Load() error {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		return fmt.Errorf("Unable to start UI: %v", err.Error())
	}

	return nil
}
