package main

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// MakeForm makes screen  Form to to insert/modify table record
func (scr *spur) MakeTopMenu(app *tview.Application) error {
	scr.topMenu = tview.NewForm()
	scr.topMenu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight || event.Key() == tcell.KeyDown {
			return tcell.NewEventKey(tcell.KeyTab, 0x09, 0)
		}
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyUp {
			return tcell.NewEventKey(tcell.KeyBacktab, 0x09, 0)
		}
		return event
	})
	scr.topMenu.SetButtonBackgroundColor(tcell.ColorDarkBlue)
	//	submit := func() {
	// 	if len(k) > 0 {
	// 		_, ok := scr.records[k]
	// 		if !ok {
	// 			scr.keys = append(scr.keys, k)
	// 		}
	// 		for j := len(v) - 1; j >= 0; j-- {
	// 			if len(v[j]) > 0 {
	// 				break
	// 			}
	// 			v = v[:j]
	// 		}
	// 		if len(v) > scr.width {
	// 			scr.width = len(v)
	// 		}
	// 		scr.records[k] = v
	// 		scr.visibility[k] = vsbl
	// 		scr.table.Clear()
	// 		scr.flex.RemoveItem(scr.table)
	// 		scr.MakeTable(app)
	// 		scr.flex.AddItem(scr.table, 0, 2, false)
	// 		scr.form.Clear(true)
	// 		scr.ChangeState(StateAlert)
	// 	}
	// scr.table.SetSelectable(true, true)
	// app.SetFocus(scr.table)
	// for ix, key := range scr.keys {
	// 	if key == k {
	// 		scr.table.Select(ix, 0)
	// 		break
	// 	}
	// }
	//	cell := scr.table.GetCell(scr.table.GetSelection())
	//	clipboard.WriteAll(cell.Text)
	//}
	// cancel := func() {
	// 	scr.table.SetSelectable(true, true)
	// 	app.SetFocus(scr.table)
	// 	for ix, key := range scr.keys {
	// 		if key == k {
	// 			scr.table.Select(ix, 0)
	// 			break
	// 		}
	// 	}
	// 	scr.form.Clear(true)
	// 	cell := scr.table.GetCell(scr.table.GetSelection())
	// 	clipboard.WriteAll(cell.Text)
	// }
	fselect := func() {
		scr.table.SetSelectable(true, true)
		app.SetFocus(scr.table)
		row, col := scr.table.GetSelection()
		if row < 1 || col < 1 {
			row = 1
			col = 1
		}
		//cell := scr.table.GetCell(row, col)
		scr.activeRow = row
		scr.activeColumn = col
		//clipboard.WriteAll(cell.Text)
	}
	scr.topMenu.AddButton("Select", fselect)

	addHidden := func() {
		scr.activeRow = -1
		scr.MakeForm(app, "h")
		modal := CompoundModal(scr.form, 45, 15)
		scr.root = scr.root.AddPage(ModalName, modal, true, true)
		app.SetRoot(scr.root, true)
		app.SetFocus(modal)
	}

	scr.topMenu.AddButton("Add", addHidden)

	scr.topMenu.AddButton("Edit", func() {
		visibility := "v"
		if scr.activeRow > 0 {
			visibility = scr.visibility[scr.keys[scr.activeRow-1]]
		}
		scr.MakeForm(app, visibility)
		modal := CompoundModal(scr.form, 45, 15)
		scr.root = scr.root.AddPage(ModalName, modal, true, true)
		app.SetRoot(scr.root, true)
		app.SetFocus(modal)

	})

	scr.topMenu.AddButton("Delete", func() {
		modal := tview.NewModal()
		var key string
		if scr.activeRow > 0 {
			key = scr.keys[scr.activeRow-1]
		}
		if len(key) > 0 {
			modal.SetText("Delete record:" + key + "?")
			modal.AddButtons([]string{"Delete", "Cancel"})
		} else {
			modal.SetText("Nothing to delete. Record empty")
			modal.AddButtons([]string{"OK"})
		}
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				scr.UpdateRecords(key, nil, "")
				scr.activeRow--
				if scr.activeRow <= 0 {
					scr.activeRow = 1
				}
				scr.activeColumn = 1
				scr.UpdateTable(app)
			}
			scr.root.RemovePage(ModalName)
			scr.table.SetSelectable(true, true)
			app.SetFocus(scr.table)
		})
		modalo := CompoundModal(modal, 15, 5)
		scr.root = scr.root.AddPage(ModalName, modalo, true, true)
		app.SetRoot(scr.root, true)
		app.SetFocus(modal)
	})

	scr.topMenu.AddButton("Save", func() {
		//scr.MakeSaveForm(app, "")
		modal := tview.NewModal()
		modal.SetText("Save " + scr.cribName + "?")
		modal.AddButtons([]string{"Save", "Cancel"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				scr.Save()
				scr.topMenu.GetButton(4).SetLabel("Save")
			}
			scr.root.RemovePage(ModalName)
			app.SetFocus(scr.topMenu)
		})
		modalo := CompoundModal(modal, 15, 5)
		scr.root = scr.root.AddPage(ModalName, modalo, true, true)
		app.SetRoot(scr.root, true)
		app.SetFocus(modal)
	})
	scr.topMenu.AddButton("Password", func() {
		scr.MakeNewPasswordForm(app)
		pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		pwdFlex.AddItem(scr.form, 0, 2, true)
		pwdFlex.SetBackgroundColor(tcell.ColorBlue)
		pwdFlex.SetBorder(true) // In case of true border is on black background
		//pwdFlex.SetBorderPadding(2, 2, 2, 2)
		modal := CompoundModal(pwdFlex, 28, 14)
		scr.root = scr.root.AddPage(ModalName, modal, true, true)
		app.SetRoot(scr.root, true)
		app.SetFocus(modal)
	})

	fexit := func() {
		clipboard.WriteAll("")
		app.Stop()
	}

	scr.topMenu.AddButton("Exit", fexit)
	return nil
}
