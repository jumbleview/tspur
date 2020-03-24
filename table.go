package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const hiddenText = " **************** "

func (scr *spur) ToClipBoard(row int, column int) {
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

// AttachData initialize spur and attaches data to it
func (scr *spur) AttachData(data []byte, pswd string) {
	scr.passwd = pswd
	var sdata []string
	scr.records = make(map[string][]string)
	scr.visibility = make(map[string]string)

	if len(data) > 0 {
		sdata = strings.Split(string(data), "\n")
	}
	for _, s := range sdata {
		// parse string as csv
		elems := strings.Split(s, ",")
		if len(elems) > 1 {
			values := elems[2:]
			scr.keys = append(scr.keys, elems[1])
			if len(values) > scr.width {
				scr.width = len(values)
			}
			scr.records[elems[1]] = values
			scr.visibility[elems[1]] = elems[0]
		}
	}
}

func (scr *spur) Hide(row int, column int) {
	if row < 1 || column < 2 {
		return // nothing ot do: hever hide key
	}
	key := scr.keys[row-1]
	if scr.visibility[key] != "h" {
		return // Row isn't hidden, nothing to do.
	}
	values := scr.records[key]
	if column < len(values)+2 {
		scr.table.GetCell(row, column).SetText(hiddenText)
	}
}

func (scr *spur) Visualize(row int, column int) {
	if row < 1 || column < 1 {
		return
	}
	key := scr.keys[row-1]
	values := scr.records[key]
	value := key
	if column > 1 {
		if column < len(values)+2 {
			value = values[column-2]
			scr.table.GetCell(row, column).SetText(value)
		}
	}
}

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
	sort.Slice(scr.keys, func(i, j int) bool {
		return strings.ToLower(scr.keys[i]) < strings.ToLower(scr.keys[j])
	})
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
	width := scr.width
	if width < 2 {
		width = 2
	}
	for c := 0; c < width; c++ {
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
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter).SetSelectable(true))
		values := scr.records[key]
		for c := 0; c < width; c++ {
			value := ""
			if c < len(values) {
				value = values[c]
			}
			tblValue := value
			if scr.visibility[key] == "h" {

				if len(tblValue) > 0 {
					tblValue = hiddenText
				} else {
					tblValue = ""
				}
			}
			scr.table.SetCell(r+1, c+2, tview.NewTableCell(tblValue).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).SetSelectable(true))
		}
	}
	return nil
}

// MakeBaseTable makes table to visualise at program start and assignes methods
func (scr *spur) MakeBaseTable(app *tview.Application) {
	scr.mode = ModeClipEnter
	scr.modes = tview.NewTable().SetBorders(false)
	scr.table = tview.NewTable().SetBorders(true)
	scr.table.SetBordersColor(tcell.ColorYellow)
	// Making table title
	scr.table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("#24")).
		SetTextColor(tcell.ColorYellow).SetSelectable(false))
	scr.table.SetCell(0, 1, tview.NewTableCell("Record Name").
		SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).
		SetSelectable(false))
	for i := 0; i < 3; i++ {
		scr.table.SetCell(0, i+2, tview.NewTableCell(fmt.Sprintf("Field %d", i)).
			SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	// Making table body
	const hight = 24
	const width = 3
	for r := 0; r < hight; r++ {
		scr.table.SetCell(r+1, 0, tview.NewTableCell(fmt.Sprintf("%d", r+1)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter).SetSelectable(false))
		for j := 1; j <= width; j++ {
			scr.table.SetCell(r+1, j, tview.NewTableCell(hiddenText).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).SetSelectable(true))
		}
	}
	scr.table.SetSelectedFunc(func(row int, column int) {
		cell := scr.table.GetCell(row, column)
		cell.SetTextColor(tcell.ColorRed)
		if len(scr.records) <= 0 {
			return
		}
		switch scr.mode {
		case ModeClipEnter:
			scr.ToClipBoard(row, column)
		case ModeClipSelect: // nothing to do
		case ModeVisibleEnter:
			scr.Visualize(row, column)
		case ModeVisibleSelect: // nothing to do
		}
		scr.activeRow = row
		scr.activeColumn = column
	})
	scr.table.SetSelectionChangedFunc(func(row int, column int) {
		if (row < 1) || (column < 1) {
			return
		}
		switch scr.mode {
		case ModeClipSelect:
			scr.ToClipBoard(row, column)
		case ModeVisibleSelect:
			scr.Visualize(row, column)
			fallthrough
		case ModeVisibleEnter:
			scr.Hide(scr.activeRow, scr.activeColumn)
		}
		scr.activeRow = row
		scr.activeColumn = column
	})
	scr.table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
		} else if key == tcell.KeyEscape {
			scr.Hide(scr.activeRow, scr.activeColumn)
			scr.table.SetSelectable(false, false)
			app.SetFocus(scr.topMenu)
		}
	})
	scr.table.SetFixed(1, 1)
}

func (scr *spur) MoveFocusToTable(app *tview.Application) {
	scr.table.SetSelectable(true, true)
	app.SetFocus(scr.table)
	scr.arrowBarrier = ArrowDefaultBarrier
}
