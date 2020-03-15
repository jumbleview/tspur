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
		if scr.mode == ModeVisualSelect {
			scr.Visualize(row, col)
		}
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
		modal.SetBackgroundColor(tcell.ColorDarkCyan)
		modal.SetButtonBackgroundColor(tcell.ColorDarkCyan)
		modal.SetTextColor(tcell.ColorWhite)
		modal.SetButtonTextColor(tcell.ColorWhite)
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
				scr.root.RemovePage(ModalName)
				scr.table.SetSelectable(true, true)
				app.SetFocus(scr.table)
				scr.topMenu.GetButton(4).SetLabel("Save!")
			} else {
				scr.root.RemovePage(ModalName)
				app.SetFocus(scr.topMenu)
			}
		})
		modalo := CompoundModal(modal, 15, 5)
		scr.root = scr.root.AddPage(ModalName, modalo, true, true)
		app.SetRoot(scr.root, true)
		app.SetFocus(modal)
	})

	scr.topMenu.AddButton("Save", func() {
		//scr.MakeSaveForm(app, "")
		modal := tview.NewModal()
		modal.SetBackgroundColor(tcell.ColorDarkCyan)
		modal.SetButtonBackgroundColor(tcell.ColorDarkCyan)
		modal.SetTextColor(tcell.ColorWhite)
		modal.SetButtonTextColor(tcell.ColorWhite)
		modal.SetText("Save page?")
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
	scr.topMenu.AddButton("Mode", func() {
		scr.MakeModeTable(app)
	})

	scr.topMenu.AddButton("Password", func() {
		needOldPassword := true
		scr.MakeNewPasswordForm(app, " Change page password ", needOldPassword)
	})

	fexit := func() {
		clipboard.WriteAll("")
		app.Stop()
	}

	scr.topMenu.AddButton("Exit", fexit)
	return nil
}
