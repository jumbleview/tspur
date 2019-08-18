package main

import (
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/rivo/tview"
)

func (scr *spur) MakePasswordForm(app *tview.Application) error {
	scr.form = tview.NewForm()
	scr.form.AddPasswordField("password", "", 12, '^', func(inp string) {
		scr.passwd = inp
	})
	return nil
}

// MakeForm makes screen  Form to to insert/modify table record
func (scr *spur) MakeForm(app *tview.Application, vsbl string) error {
	scr.form = tview.NewForm()
	count := scr.width
	var k string
	var v []string
	if len(scr.keys) > 0 {
		if scr.activeRow >= 0 {
			k = scr.keys[scr.activeRow]
		}
		scr.form.AddInputField("Key", k, 21, nil, func(inp string) {
			k = inp
		})
		if len(k) > 0 {
			v = append(v, scr.records[k]...)
		}
		for i := 0; i <= count; i++ {
			valName := strconv.Itoa(i)
			if i == count {
				valName = "+"
			}
			if i >= len(v) {
				v = append(v, "")
			}
			locali := i
			accepted := func(inp string, last rune) bool {
				clipboard.WriteAll(inp)
				return true
			}
			changed := func(inp string) {
				v[locali] = inp
				clipboard.WriteAll(inp)
			}
			if vsbl == "" {
				vsbl = scr.visibility[k]
			}
			if vsbl == "h" {
				scr.form.AddPasswordField(valName, v[i], 21, '*', changed)
			} else {
				scr.form.AddInputField(valName, v[i], 21, accepted, changed)
			}
		}
	} else {
		v = append(v, "")
		scr.form.AddInputField("Key", "", 21, nil, func(inp string) {
			k = inp
		})
		scr.form.AddInputField("+", "", 21, nil, func(inp string) {
			v[0] = inp
		})
	}
	submit := func() {
		if len(k) > 0 {
			_, ok := scr.records[k]
			if !ok {
				scr.keys = append(scr.keys, k)
			}
			for j := len(v) - 1; j >= 0; j-- {
				if len(v[j]) > 0 {
					break
				}
				v = v[:j]
			}
			if len(v) > scr.width {
				scr.width = len(v)
			}
			scr.records[k] = v
			scr.visibility[k] = vsbl
			scr.table.Clear()
			scr.flex.RemoveItem(scr.table)
			scr.MakeTable(app)
			scr.flex.AddItem(scr.table, 0, 2, false)
			scr.form.Clear(true)
			scr.ChangeState(StateAlert)
		}
		scr.table.SetSelectable(true, true)
		app.SetFocus(scr.table)
		for ix, key := range scr.keys {
			if key == k {
				scr.table.Select(ix, 0)
				break
			}
		}
		cell := scr.table.GetCell(scr.table.GetSelection())
		clipboard.WriteAll(cell.Text)
	}
	cancel := func() {
		scr.table.SetSelectable(true, true)
		app.SetFocus(scr.table)
		for ix, key := range scr.keys {
			if key == k {
				scr.table.Select(ix, 0)
				break
			}
		}
		scr.form.Clear(true)
		cell := scr.table.GetCell(scr.table.GetSelection())
		clipboard.WriteAll(cell.Text)
	}
	//scr.form.AddButton("Set", submit)
	//scr.form.AddButton("Esc", cancel)
	scr.form.AddButton("Submit", func() {
		ttl := "Store " + k + "?"
		modal := tview.NewModal().SetText(ttl)
		modal.AddButtons([]string{"Yes", "No"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				submit()
			} else {
				cancel()
			}
			scr.flex.RemoveItem(modal)
			//app.SetFocus(scr.flex)
		})
		scr.flex.AddItem(modal, 20, 1, true)
		app.SetFocus(modal)
	})

	scr.form.SetCancelFunc(cancel)

	return nil
}
