package main

import (
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// MakeTable makes table out of parsed data
func (scr *spur) MakeTable(app *tview.Application) error {
	sort.Strings(scr.keys)
	scr.width = 0
	for _, v := range scr.records {
		if len(v) > scr.width {
			scr.width = len(v)
		}
	}
	scr.table = tview.NewTable().SetBorders(true)
	rowsMax := len(scr.keys)
	for r := 0; r < rowsMax; r++ {
		key := scr.keys[r]
		color := tcell.ColorWhite
		c := 0
		scr.table.SetCell(r, c,
			tview.NewTableCell(key).
				SetTextColor(color).
				SetAlign(tview.AlignCenter))
		values := scr.records[key]
		for c := 0; c < scr.width; c++ {
			value := ""
			if c < len(values) {
				value = values[c]
			}
			tblValue := value
			if scr.visibility[key] == "h" {
				tblValue = strings.Repeat("*", len(value))
			}
			scr.table.SetCell(r, c+1,
				tview.NewTableCell(tblValue).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
		}
	}
	toClipBoard := func(row int, column int) {
		key := scr.keys[row]
		values := scr.records[key]
		value := key
		if column > 0 {
			if column <= len(values) {
				value = values[column-1]
			} else {
				value = ""
			}
		}
		clipboard.WriteAll(value)
	}
	scr.table.SetSelectedFunc(func(row int, column int) {
		cell := scr.table.GetCell(row, column)
		cell.SetTextColor(tcell.ColorRed)
		if len(scr.records) > 0 {
			toClipBoard(row, column)
		}
		scr.table.SetSelectable(false, false)
		scr.activeRow = row
		app.SetFocus(scr.list)
	})
	scr.table.SetSelectionChangedFunc(func(row int, column int) {
		//cell := scr.table.GetCell(row, column)
		//clipboard.WriteAll(cell.Text)
		if len(scr.records) > 0 {
			toClipBoard(row, column)
		}
		scr.activeRow = row
		//table.SetSelectable(false, false)
	})
	//scr.table.SetBorder(true).SetTitle("spur-table")
	//scr.table.SetTitle("spur-table")
	scr.table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			//table.SetSelectable(true, true)
			scr.table.SetSelectable(false, false)
			app.SetFocus(scr.list)
		}
		if key == tcell.KeyEscape {
			scr.table.SetSelectable(false, false)
			app.SetFocus(scr.list)
		}
	})
	scr.table.SetFixed(0, 1)

	return nil
}

func (scr *spur) MakeChangeState() error {
	scr.changeState = tview.NewTable().SetBorders(false)
	scr.changeState.SetCell(0, 0,
		tview.NewTableCell(StateSaved).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter))
	scr.changeState.SetSelectable(false, false)
	return nil
}

func (scr *spur) ChangeState(state string) {
	scr.changeState.GetCell(0, 0).SetText(state)
}

func (scr *spur) IsStateAlter() bool {
	text := scr.changeState.GetCell(0, 0).Text
	return (text != StateSaved)
}
