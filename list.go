package main

import (
	"github.com/atotto/clipboard"
	"github.com/rivo/tview"
)

// MakeList makes screen list to with list of application action
func (scr *spur) MakeList(app *tview.Application) error {
	scr.list = tview.NewList()
	scr.list.ShowSecondaryText(false)
	scr.list.SetSelectedFocusOnly(true)
	lnavigate := func() {
		scr.table.SetBorders(true).SetSelectable(true, true)
		app.SetFocus(scr.table)
		cell := scr.table.GetCell(scr.table.GetSelection())
		clipboard.WriteAll(cell.Text)
	}
	modify := func() {
		scr.form.Clear(true)
		scr.lstFlx.RemoveItem(scr.form)
		//scr.lstFlx.RemoveItem(scr.list)
		scr.MakeForm(app, "")
		//scr.lstFlx.AddItem(scr.list, 0, 1, true)
		scr.lstFlx.AddItem(scr.form, 0, 2, false)
		app.SetFocus(scr.form)
	}
	create := func() {
		scr.form.Clear(true)
		scr.lstFlx.RemoveItem(scr.form)
		//scr.lstFlx.RemoveItem(scr.list)
		scr.activeRow = -1
		scr.MakeForm(app, "v")
		scr.activeRow = 0
		//scr.lstFlx.AddItem(scr.list, 0, 1, true)
		scr.lstFlx.AddItem(scr.form, 0, 2, false)
		app.SetFocus(scr.form)
	}

	delete := func() {
		if len(scr.keys) > 0 {
			if scr.activeRow >= 0 {
				k := scr.keys[scr.activeRow]
				ttl := "Delete " + k + "?"
				modal := tview.NewModal().SetText(ttl)
				modal.AddButtons([]string{"Yes", "No"})
				modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						_, ok := scr.records[k]
						if ok {
							for i := range scr.keys {
								if scr.keys[i] == k {
									// Found!
									scr.keys = append(scr.keys[:i], scr.keys[i+1:]...)
									delete(scr.records, k)
									scr.ChangeState(StateAlter)
									break
								}
							}
						}
					}
					scr.flex.RemoveItem(modal)
					scr.table.Clear()
					scr.flex.RemoveItem(scr.table)
					scr.MakeTable(app)
					scr.flex.AddItem(scr.table, 0, 2, false)
					scr.table.SetSelectable(false, false)
					app.SetFocus(scr.flex)
				})
				scr.flex.AddItem(modal, 20, 1, true)
				app.SetFocus(modal)
			}
		}
	}

	hidden := func() {
		scr.form.Clear(true)
		scr.lstFlx.RemoveItem(scr.form)
		//scr.lstFlx.RemoveItem(scr.list)
		scr.activeRow = -1
		scr.MakeForm(app, "h")
		scr.activeRow = 0
		//scr.lstFlx.AddItem(scr.list, 0, 1, true)
		scr.lstFlx.AddItem(scr.form, 0, 2, false)
		app.SetFocus(scr.form)
	}

	save := func() {
		csv := ""
		for _, key := range scr.keys {
			line := scr.visibility[key]
			line += ","
			line += key
			values := scr.records[key]
			for _, value := range values {
				line += ","
				line += value
			}
			line += "\n"
			csv += line
		}
		if len(csv) > 0 {
			//crib.Write([]byte(csv))
			err := EncryptFile(CribName, []byte(csv), scr.passwd)
			if err != nil {
				panic(err.Error())
			}
			scr.ChangeState(StateSaved)
		}
		//app.Stop()
		lnavigate()
	}

	lsave := func() {
		modal := tview.NewModal().SetText("Save, Really?")
		modal.AddButtons([]string{"Yes", "No"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				save()
			}
			scr.flex.RemoveItem(modal)
			scr.table.SetSelectable(false, false)
			app.SetFocus(scr.flex)
		})
		scr.flex.AddItem(modal, 20, 1, true)
		app.SetFocus(modal)
	}

	lexit := func() {
		if scr.IsStateAlter() {
			modal := tview.NewModal().SetText("Table changed. Exit? Really?")
			modal.AddButtons([]string{"No", "Yes"})
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					app.Stop()
				} else {
					scr.flex.RemoveItem(modal)
					app.SetFocus(scr.flex)
				}
			})
			scr.flex.AddItem(modal, 20, 1, true)
			app.SetFocus(modal)
		} else {
			app.Stop()
		}
	}
	scr.list.AddItem("Navigate", "", 'n', lnavigate)
	scr.list.AddItem("Modify", "", 'm', modify)
	scr.list.AddItem("+Visible", "", 'v', create)
	scr.list.AddItem("+Hidden", "", 'h', hidden)
	scr.list.AddItem("Delete", "", 'd', delete)
	scr.list.AddItem("Save", "", 's', lsave)
	scr.list.AddItem("Exit", "", 'e', lexit)
	return nil
}
