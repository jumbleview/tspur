package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// UpdateRecords add or update element to the records or delete it if values=nil
func (scr *spur) UpdateRecords(key string, values []string, visibility string) int {
	if values == nil { // delete element from the records
		delete(scr.records, key)
		delete(scr.visibility, key)
	} else { // add or update records and visibility maps
		scr.records[key] = values
		scr.visibility[key] = visibility
	}
	// rebuild and sort slice with keys and update max width
	scr.keys = scr.keys[:0] // empties keys slice
	scr.width = 0
	for k, v := range scr.records {
		scr.keys = append(scr.keys, k) // rebuild keys slice
		if scr.width < len(v) {
			scr.width = len(v)
		}
	}
	sort.Strings(scr.keys)
	i := 0
	k := ""
	for i, k = range scr.keys {
		if k == key {
			break
		}
	}
	return i
}

// UpdateTable brings table in accordance with records
func (scr *spur) UpdateTable(app *tview.Application) error {
	scr.table.Clear().SetBorders(true)
	scr.table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("#%d", len(scr.records))).
		SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).
		SetSelectable(false))
	scr.table.SetCell(0, 1, tview.NewTableCell("Record Name").
		SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).
		SetSelectable(false))
	for c := 0; c < scr.width; c++ {
		scr.table.SetCell(0, c+2, tview.NewTableCell(fmt.Sprintf("Field %d", c)).
			SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	for r := 0; r < len(scr.keys); r++ {
		key := scr.keys[r]
		scr.table.SetCell(r+1, 0, tview.NewTableCell(fmt.Sprintf("%d", r+1)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter).SetSelectable(false))
		scr.table.SetCell(r+1, 1, tview.NewTableCell(key).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).SetSelectable(true))
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
			scr.table.SetCell(r+1, c+2, tview.NewTableCell(tblValue).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).SetSelectable(true))
		}
	}
	return nil
}

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
	scr.table.SetBordersColor(tcell.ColorYellow)
	err := scr.UpdateTable(app)
	if err != nil {
		return err
	}
	toClipBoard := func(row int, column int) {
		if row < 1 || column < 1 {
			return
		}
		key := scr.keys[row-1]
		values := scr.records[key]
		value := key
		if column > 1 {
			if column < len(values)+2 {
				value = values[column-2]
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
		//scr.table.SetSelectable(false, false)
		scr.activeRow = row
		scr.activeColumn = column
		//app.SetFocus(scr.topMenu)
	})
	scr.table.SetSelectionChangedFunc(func(row int, column int) {
		//if len(scr.records) > 0 {
		if (row < 1) || (column < 1) {
			return
		}
		//toClipBoard(row, column)
		scr.activeRow = row
		scr.activeColumn = column
		//table.SetSelectable(false, false)
	})
	scr.table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			// scr.table.SetSelectable(true, true)
			// app.SetFocus(scr.table)
			// //scr.table.SetSelectable(false, false)
			// //app.SetFocus(scr.topMenu)
			// toClipBoard(scr.activeRow, scr.activeColumn)
		} else if key == tcell.KeyEscape {
			scr.table.SetSelectable(false, false)
			app.SetFocus(scr.topMenu)
		}
	})
	scr.table.SetFixed(1, 1)

	return nil
}
