package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/google/logger"
	"github.com/rivo/tview"
)

const (
	cols = "ABCDEFGHI"
	rows = "123456789"
)

var (
	events    []string
	app       *tview.Application
	attempted map[string]bool
	hits      int
)

// Load generates and displays the UI to the user.
func Load(send, recv chan string, closed chan bool) error {
	go func() {
		<-closed // signals that the server closed the connection
		app.Stop()
		logger.Errorf("The server closed the connection")
	}()
	attempted = make(map[string]bool)
	app = tview.NewApplication()
	pages := tview.NewPages()

	title := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Gottleships")
	sbTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Event Log")
	sbText := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	mainCol := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 3, 1, false)
	sidebar := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(sbTitle, 3, 1, false).
		AddItem(sbText, 0, 1, false)

	grid := tview.NewGrid().
		SetColumns(0, 25).
		SetBorders(true)
	grid.AddItem(mainCol, 0, 0, 1, 2, 0, 0, true).
		AddItem(mainCol, 0, 0, 1, 1, 0, 100, true).
		// Sidebar won't show if window is less than 100 cells wide
		AddItem(sidebar, 0, 0, 0, 0, 0, 0, false).
		AddItem(sidebar, 0, 1, 1, 1, 0, 100, false)

	pages.AddPage("main", grid, true, true)

	board := setupBoard()
	// When a cell is selected in the TUI, a handler func is called for that event.
	// This function sends and receives over the channel passed to the function from the client.
	// We check the response from the server, updates the UI elements based on hit or miss, and
	// displays a game over modal when all ships have been sunk.
	board.SetSelectedFunc(func(row, col int) {
		cell := string(cols[col-1]) + string(rows[row-1])
		// Check if already attempted
		if _, ok := attempted[cell]; ok {
			duplicate := tview.NewModal().
				SetText(fmt.Sprintf("You've already fired on that target!")).
				AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.HidePage("duplicate")
			})
			pages.AddPage("duplicate", duplicate, false, true)
			return
		}
		send <- cell

		events = append([]string{fmt.Sprintf("[yellow]SENT: [white]%v", cell)}, events...)
		sbText.SetText(strings.Join(events, "\n"))

		resp, ok := <-recv
		if !ok {
			app.Stop()
		}
		switch resp {
		case "HIT":
			events = append([]string{fmt.Sprintf("[yellow]RECV: [green]%v", resp)}, events...)
			sbText.SetText(strings.Join(events, "\n"))
			board.SetCell(row, col, tview.NewTableCell("[red]X").
				SetAlign(tview.AlignCenter))
			hits++
			attempted[cell] = true
		case "MISS":
			events = append([]string{fmt.Sprintf("[yellow]RECV: [red]%v", resp)}, events...)
			sbText.SetText(strings.Join(events, "\n")).ScrollToBeginning()
			board.SetCell(row, col, tview.NewTableCell("[grey]--").
				SetAlign(tview.AlignCenter))
			attempted[cell] = false
		}

		if hits == 14 {
			resp, ok := <-recv
			if !ok {
				app.Stop()
			}
			i, err := strconv.Atoi(resp)
			// 9x9 grid is equal to 81 possible shots
			if err != nil || i > 81 {
				logger.Fatal("Received an invalid response from the serveur, existing...")
			}
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Game Over! Score: [yellow]%v\n[white]Press Enter to Exit", resp)).
				AddButtons([]string{"Exit"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.Stop()
			})
			pages.AddPage("Game Over", modal, false, true)
		}
	})

	boardFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(board, 0, 1, true)

	// The gaps to the sides of the board will only show if the
	// width of the container is greater than 80. This stops the board getting
	// too wide.
	innerGrid := tview.NewGrid().
		SetColumns(0, -4, 0).
		AddItem(boardFlex, 0, 0, 1, 3, 0, 0, true).
		AddItem(boardFlex, 0, 1, 1, 1, 0, 80, true).
		AddItem(tview.NewTextView(), 0, 0, 1, 1, 0, 80, false).
		AddItem(tview.NewTextView(), 0, 2, 1, 1, 0, 80, false)
	mainCol.AddItem(innerGrid, 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		return event
	})

	quit := tview.NewTextView().
		SetText("Ctrl+Q: Quit  |  ArrowKeys: Move  |  Enter: Fire!")
	main := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(quit, 1, 1, false)

	if err := app.SetRoot(main, true).Run(); err != nil {
		return fmt.Errorf("Unable to start UI: %v", err.Error())
	}

	return nil
}

// setupBoard returns a table with that represents the battleship
// game board. The axis are labels, and the inner cells are all selectable.
func setupBoard() *tview.Table {
	board := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, true).
		SetBorders(true).
		SetCell(0, 0, tview.NewTableCell("").
			SetSelectable(false).SetExpansion(1))

	for pos, char := range cols {
		col := tview.NewTableCell(string(char)).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(1).
			SetTextColor(tcell.ColorAqua)
		row := tview.NewTableCell(string(rows[pos])).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(1).
			SetTextColor(tcell.ColorYellow)

		board.SetCell(pos+1, 0, row)
		board.SetCell(0, pos+1, col)

		for i := range rows {
			item := tview.NewTableCell(" ● ").
				SetSelectable(true).
				SetExpansion(1).
				SetAlign(tview.AlignCenter)
			board.SetCell(pos+1, i+1, item)
		}
	}

	return board
}
