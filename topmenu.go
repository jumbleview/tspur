package main

import (
	"fmt"

	//"github.com/atotto/clipboard"
	"github.com/d-tsuji/clipboard"
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

	spr.topMenu.AddButton("Select", func() {
		spr.MoveFocusToTable(app)
		row, col := spr.table.GetSelection()
		if row < 1 || col < 1 {
			row = 1
			col = 1
		}
		spr.activeRow = row
		spr.activeColumn = col
		if spr.mode == ModeVisibleSelect {
			spr.Visualize(row, col)
		} else if spr.mode == ModeClipSelect {
			spr.ToClipBoard(row, col)
		}
		spr.arrowBarrier = -1
	})

	spr.topMenu.AddButton("Mode", func() {
		spr.MakeModeTable(app)
	})

	spr.topMenu.AddButton("Add", func() {
		spr.activeRow = -1
		spr.MakeForm(app, "h")
		modal := CompoundModal(spr.form, 45, 15)
		spr.root = spr.root.AddPage(ModalName, modal, true, true)
		app.SetRoot(spr.root, true)
		app.SetFocus(modal)
	})

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
		modal := spr.MakeNewModal()
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
		modal := spr.MakeNewModal()
		modal.SetText("Save page?")
		modal.AddButtons([]string{"Save", "Cancel"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				spr.Save()
				spr.topMenu.GetButton(spr.saveMenuInx).SetLabel("Save")
				if spr.topMenu.GetButton(spr.saveMenuInx+1).GetLabel() == "Git" {
					spr.topMenu.GetButton(spr.saveMenuInx + 1).SetLabel("Git!")
				}
			}
			spr.root.RemovePage(ModalName)
			app.SetFocus(spr.topMenu)
		})
		modalo := CompoundModal(modal, 15, 5)
		spr.root = spr.root.AddPage(ModalName, modalo, true, true)
		app.SetRoot(spr.root, true)
		app.SetFocus(modal)
	})
	if CheckGit(spr.cribPath) == nil {
		spr.topMenu.AddButton("Git", func() {
			modal := spr.MakeNewModal()
			modal.SetText("Push to remote?")
			modal.AddButtons([]string{"Yes", "No"})
			attempt := 0
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if attempt == 0 && buttonLabel == "Yes" {
					modal.SetText("Pushing...")
					app.ForceDraw()
					txt, err := PushRemote(spr.cribPath, spr.cribBase, spr.commits)
					spr.commits = spr.commits[:0] // make it empty for the next commit
					attempt = 1
					if err != nil {
						txt = fmt.Sprintf("%s%s", txt, err)
					}
					modal.SetText(txt)
					modal.ClearButtons()
					modal.AddButtons([]string{"OK"})
					app.ForceDraw()
				} else {
					spr.topMenu.GetButton(spr.saveMenuInx + 1).SetLabel("Git")
					spr.root.RemovePage(ModalName)
					app.SetFocus(spr.topMenu)
				}
			})
			modalo := CompoundModal(modal, 15, 5)
			spr.root = spr.root.AddPage(ModalName, modalo, true, true)
			app.SetRoot(spr.root, true)
			app.SetFocus(modal)
		})

	}
	spr.topMenu.AddButton("Password", func() {
		needOldPassword := true
		spr.MakeNewPasswordForm(app, " Change page password ", needOldPassword)
	})

	spr.topMenu.AddButton("Exit", func() {
		saveLabel := spr.topMenu.GetButton(spr.saveMenuInx).GetLabel()
		if saveLabel == "Save" { // nothing to save. Just exit
			//clipboard.WriteAll("")
			clipboard.Set("")
			app.Stop()
			return
		}
		modal := spr.MakeNewModal()
		modal.SetText("Page not saved. Exit?")
		modal.AddButtons([]string{"Exit", "Cancel"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Exit" {
				//clipboard.WriteAll("")
				clipboard.Set("")
				app.Stop()
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

	spr.saveMenuInx = spr.topMenu.GetButtonIndex("Save")
	return nil
}
