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

// Format used to make table title
const FmtFieldTitle = "     Field %d     "

func (spr *Spur) ToClipBoard(row int, column int) {
	if row < 1 || column < 1 {
		return
	}
	key := spr.keys[row-1]
	values := spr.records[key]
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
func (spr *Spur) AttachData(data []byte, pswd string) {
	spr.passwd = pswd
	var sdata []string
	spr.records = make(map[string][]string)
	spr.visibility = make(map[string]string)

	if len(data) > 0 {
		sdata = strings.Split(string(data), "\n")
	}
	for _, s := range sdata {
		// parse string as csv
		elems := strings.Split(s, ",")
		if len(elems) > 1 {
			values := elems[2:]
			spr.keys = append(spr.keys, elems[1])
			if len(values) > spr.width {
				spr.width = len(values)
			}
			spr.records[elems[1]] = values
			spr.visibility[elems[1]] = elems[0]
		}
	}
}

func (spr *Spur) Hide(row int, column int) {
	if row < 1 || column < 2 {
		return // nothing ot do: hever hide key
	}
	key := spr.keys[row-1]
	if spr.visibility[key] != "h" {
		return // Row isn't hidden, nothing to do.
	}
	values := spr.records[key]
	if column < len(values)+2 {
		spr.table.GetCell(row, column).SetText(hiddenText)
	}
}

func (spr *Spur) Visualize(row int, column int) {
	if row < 1 || column < 1 {
		return
	}
	key := spr.keys[row-1]
	values := spr.records[key]
	value := key
	if column > 1 {
		if column < len(values)+2 {
			value = values[column-2]
			spr.table.GetCell(row, column).SetText(value)
		}
	}
}

// UpdateRecords add or update element to the records or delete it if values=nil
func (spr *Spur) UpdateRecords(key string, values []string, visibility string) int {
	if values == nil { // delete element from the records
		s := key + " deleted"
		spr.commits = append(spr.commits, s)
		delete(spr.records, key)
		delete(spr.visibility, key)
	} else { // add or update records and visibility maps
		s := key + " changed"
		spr.commits = append(spr.commits, s)
		spr.records[key] = values
		spr.visibility[key] = visibility
	}
	// rebuild and sort slice with keys and update max width
	spr.keys = spr.keys[:0] // empties keys slice
	spr.width = 0
	for k, v := range spr.records {
		spr.keys = append(spr.keys, k) // rebuild keys slice
		if spr.width < len(v) {
			spr.width = len(v)
		}
	}
	sort.Slice(spr.keys, func(i, j int) bool {
		return strings.ToLower(spr.keys[i]) < strings.ToLower(spr.keys[j])
	})
	i := 0
	k := ""
	for i, k = range spr.keys {
		if k == key {
			break
		}
	}
	return i
}

// UpdateTable brings table in accordance with records
func (spr *Spur) UpdateTable(app *tview.Application) error {
	spr.table.Clear().SetBorders(true)
	spr.table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("#%d", len(spr.records))).
		SetTextColor(spr.AccentColor).SetAlign(tview.AlignCenter).
		SetSelectable(false))
	spr.table.SetCell(0, 1, tview.NewTableCell("Record Name").
		SetTextColor(spr.AccentColor).SetAlign(tview.AlignCenter).
		SetSelectable(false))
	width := spr.width
	if width < 2 {
		width = 2
	}
	wmax := width
	if wmax < 3 {
		wmax = 3
	}
	for c := 0; c < 3; c++ {
		spr.table.SetCell(0, c+2, tview.NewTableCell(fmt.Sprintf(FmtFieldTitle, c)).
			SetTextColor(spr.AccentColor).SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	for r := 0; r < len(spr.keys); r++ {
		key := spr.keys[r]
		spr.table.SetCell(r+1, 0, tview.NewTableCell(fmt.Sprintf("%d", r+1)).
			SetTextColor(spr.MainColor).
			SetAlign(tview.AlignCenter).SetSelectable(false))
		spr.table.SetCell(r+1, 1, tview.NewTableCell(key).
			SetTextColor(spr.MainColor).
			SetAlign(tview.AlignCenter).SetSelectable(true))
		values := spr.records[key]
		for c := 0; c < width; c++ {
			value := ""
			if c < len(values) {
				value = values[c]
			}
			tblValue := value
			if spr.visibility[key] == "h" {

				if len(tblValue) > 0 {
					tblValue = hiddenText
				} else {
					tblValue = ""
				}
			}
			spr.table.SetCell(r+1, c+2, tview.NewTableCell(tblValue).
				SetTextColor(spr.MainColor).
				SetAlign(tview.AlignCenter).SetSelectable(true))
		}
	}
	return nil
}

// MakeBaseTable makes table to visualize at program start and assigns methods
func (spr *Spur) MakeBaseTable(app *tview.Application) {
	spr.mode = ModeClipEnter
	spr.modes = tview.NewTable().SetBorders(false)
	spr.table = tview.NewTable().SetBorders(true)
	spr.table.SetBackgroundColor(spr.MainBackgroundColor)
	spr.table.SetBordersColor(spr.AccentColor)
	// Making table title
	spr.table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("#24")).
		SetTextColor(spr.AccentColor).SetSelectable(false))
	spr.table.SetCell(0, 1, tview.NewTableCell("Record Name").
		SetTextColor(spr.AccentColor).SetAlign(tview.AlignCenter).
		SetSelectable(false))
	for i := 0; i < 3; i++ {
		spr.table.SetCell(0, i+2, tview.NewTableCell(fmt.Sprintf(FmtFieldTitle, i)).
			SetTextColor(spr.AccentColor).SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}
	// Making table body
	const hight = 24
	const width = 3
	for r := 0; r < hight; r++ {
		spr.table.SetCell(r+1, 0, tview.NewTableCell(fmt.Sprintf("%d", r+1)).
			SetTextColor(spr.MainColor).
			SetAlign(tview.AlignCenter).SetSelectable(false))
		for j := 1; j <= width; j++ {
			spr.table.SetCell(r+1, j, tview.NewTableCell(hiddenText).
				SetTextColor(spr.MainColor).
				SetAlign(tview.AlignCenter).SetSelectable(true))
		}
	}
	spr.table.SetSelectedFunc(func(row int, column int) {
		cell := spr.table.GetCell(row, column)
		cell.SetTextColor(spr.TrackingColor)
		if len(spr.records) <= 0 {
			return
		}
		switch spr.mode {
		case ModeClipEnter:
			spr.ToClipBoard(row, column)
		case ModeClipSelect: // nothing to do
		case ModeVisibleEnter:
			spr.Visualize(row, column)
		case ModeVisibleSelect: // nothing to do
		}
		spr.activeRow = row
		spr.activeColumn = column
	})
	spr.table.SetSelectionChangedFunc(func(row int, column int) {
		if (row < 1) || (column < 1) {
			return
		}
		switch spr.mode {
		case ModeClipSelect:
			spr.ToClipBoard(row, column)
		case ModeVisibleSelect:
			spr.Visualize(row, column)
			fallthrough
		case ModeVisibleEnter:
			spr.Hide(spr.activeRow, spr.activeColumn)
		}
		spr.activeRow = row
		spr.activeColumn = column
	})
	spr.table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
		} else if key == tcell.KeyEscape {
			spr.Hide(spr.activeRow, spr.activeColumn)
			spr.table.SetSelectable(false, false)
			app.SetFocus(spr.topMenu)
		}
	})
	spr.table.SetFixed(1, 1)
}

func (spr *Spur) MoveFocusToTable(app *tview.Application) {
	spr.table.SetSelectable(true, true)
	app.SetFocus(spr.table)
	spr.arrowBarrier = ArrowDefaultBarrier
}
