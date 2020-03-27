package main

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// MakeTopMenu makes application top menu to navigate/manipulate table
func (spr *Spur) MakeTopMenu(app *tview.Application) error {
	spr.topMenu = tview.NewForm()
	spr.topMenu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight || event.Key() == tcell.KeyDown {
			return tcell.NewEventKey(tcell.KeyTab, 0x09, 0)
		}
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyUp {
			return tcell.NewEventKey(tcell.KeyBacktab, 0x09, 0)
		}
		return event
	})

	spr.topMenu.SetBackgroundColor(spr.MainBackgroundColor)
	spr.topMenu.SetButtonBackgroundColor(spr.MainBackgroundColor)
	spr.topMenu.SetButtonTextColor(spr.AccentColor)
	fselect := func() {

		spr.MoveFocusToTable(app)

		row, col := spr.table.GetSelection()
		if row < 1 || col < 1 {
			row = 1
			col = 1
		}
		//cell := spr.table.GetCell(row, col)
		spr.activeRow = row
		spr.activeColumn = col
		if spr.mode == ModeVisibleSelect {
			spr.Visualize(row, col)
		} else if spr.mode == ModeClipSelect {
			spr.ToClipBoard(row, col)
		}
		spr.arrowBarrier = -1
	}
	spr.topMenu.AddButton("Select", fselect)

	spr.topMenu.AddButton("Mode", func() {
		spr.MakeModeTable(app)
	})

	addHidden := func() {
		spr.activeRow = -1
		spr.MakeForm(app, "h")
		modal := CompoundModal(spr.form, 45, 15)
		spr.root = spr.root.AddPage(ModalName, modal, true, true)
		app.SetRoot(spr.root, true)
		app.SetFocus(modal)
	}

	spr.topMenu.AddButton("Add", addHidden)

	spr.topMenu.AddButton("Edit", func() {
		visibility := "v"
		if spr.activeRow > 0 {
			visibility = spr.visibility[spr.keys[spr.activeRow-1]]
		}
		spr.MakeForm(app, visibility)
		modal := CompoundModal(spr.form, 45, 15)
		spr.root = spr.root.AddPage(ModalName, modal, true, true)
		app.SetRoot(spr.root, true)
		app.SetFocus(modal)

	})

	spr.topMenu.AddButton("Delete", func() {
		modal := tview.NewModal()
		modal.SetBackgroundColor(spr.FormBackgroundColor)
		modal.SetButtonBackgroundColor(spr.FormBackgroundColor)
		modal.SetTextColor(spr.FormColor)
		modal.SetButtonTextColor(spr.FormColor)
		var key string
		if spr.activeRow > 0 {
			key = spr.keys[spr.activeRow-1]
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
				spr.UpdateRecords(key, nil, "")
				spr.activeRow--
				if spr.activeRow <= 0 {
					spr.activeRow = 1
				}
				spr.activeColumn = 1
				spr.UpdateTable(app)
				spr.root.RemovePage(ModalName)
				spr.MoveFocusToTable(app)
				spr.topMenu.GetButton(spr.saveMenuInx).SetLabel("Save!")
			} else {
				spr.root.RemovePage(ModalName)
				app.SetFocus(spr.topMenu)
			}
		})
		modalo := CompoundModal(modal, 15, 5)
		spr.root = spr.root.AddPage(ModalName, modalo, true, true)
		app.SetRoot(spr.root, true)
		app.SetFocus(modal)
	})

	spr.topMenu.AddButton("Save", func() {
		//spr.MakeSaveForm(app, "")
		modal := tview.NewModal()
		modal.SetBackgroundColor(spr.FormBackgroundColor)
		modal.SetButtonBackgroundColor(spr.FormBackgroundColor)
		modal.SetTextColor(spr.FormColor)
		modal.SetButtonTextColor(spr.FormColor)
		modal.SetText("Save page?")
		modal.AddButtons([]string{"Save", "Cancel"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				spr.Save()
				spr.topMenu.GetButton(spr.saveMenuInx).SetLabel("Save")
			}
			spr.root.RemovePage(ModalName)
			app.SetFocus(spr.topMenu)
		})
		modalo := CompoundModal(modal, 15, 5)
		spr.root = spr.root.AddPage(ModalName, modalo, true, true)
		app.SetRoot(spr.root, true)
		app.SetFocus(modal)
	})

	spr.topMenu.AddButton("Password", func() {
		needOldPassword := true
		spr.MakeNewPasswordForm(app, " Change page password ", needOldPassword)
	})

	fexit := func() {
		clipboard.WriteAll("")
		app.Stop()
	}

	spr.topMenu.AddButton("Exit", fexit)
	spr.saveMenuInx = spr.topMenu.GetButtonIndex("Save")
	return nil
}
