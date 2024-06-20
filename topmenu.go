package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
)

func (spr *Spur) SelectTable(app *tview.Application) {
	spr.MoveFocusToTable(app)
	row, col := spr.table.GetSelection()
	// if row < 1 || col < 1 {
	// 	row = 1
	// 	col = 1
	// }
	spr.activeRow = row
	spr.activeColumn = col
	if spr.mode == ModeVisibleSelect {
		spr.Visualize(row, col)
	} else if spr.mode == ModeClipSelect {
		spr.ToClipBoard(row, col)
	}
	// spr.arrowBarrier = -1
}

func FirstToUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	start := s[0:1]
	rest := s[1:]
	return strings.ToUpper(start) + rest
}

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

	spr.topMenu.AddButton("Select", func() {
		spr.SelectTable(app)
	})

	spr.topMenu.SetBackgroundColor(spr.MainBackgroundColor)
	spr.topMenu.SetButtonBackgroundColor(spr.MainBackgroundColor)
	spr.topMenu.SetButtonTextColor(spr.AccentColor)
	label := FirstToUpper(spr.modeSet[spr.modeIndex])

	spr.topMenu.AddButton(label, func() {
		spr.Hide(spr.activeRow, spr.activeColumn)
		spr.table.SetSelectable(false, false)
		spr.isLastEventMouse = false
		spr.MakeModeTable(app)
	})

	spr.topMenu.AddButton("WWW", func() {
		spr.Hide(spr.activeRow, spr.activeColumn)
		spr.table.SetSelectable(false, false)

		modal := spr.MakeNewModal()
		var key, url string
		if spr.activeRow > 0 {
			key = spr.keys[spr.activeRow-1]
		}
		urlIndex := spr.activeColumn - 2
		if urlIndex >= 0 && urlIndex < len(spr.records[key]) {
			url = spr.records[key][urlIndex]
		}

		if len(url) > 0 {
			modal.SetText("Go to URL:\n[" + url + "] ?")
			modal.AddButtons([]string{"Yes", "No"})
		} else {
			modal.SetText("Not valid URL\n[" + url + "]")
			modal.AddButtons([]string{"OK"})
		}
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				if !strings.HasPrefix(url, "http://") &&
					!strings.HasPrefix(url, "https://") {
					url = "http://" + url
				}
				browser.OpenURL(url)
				spr.root.RemovePage(ModalName)
				spr.MoveFocusToTable(app)
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
	/*
		spr.topMenu.AddButton("Add", func() {
			spr.Hide(spr.activeRow, spr.activeColumn)
			spr.table.SetSelectable(false, false)

			spr.activeRow = -1
			spr.MakeForm(app, "h")
			modal := CompoundModal(spr.form, 45, 19)
			spr.root = spr.root.AddPage(ModalName, modal, true, true)
			app.SetRoot(spr.root, true)
			app.SetFocus(modal)
		})
	*/
	spr.topMenu.AddButton("Edit", func() {
		spr.Hide(spr.activeRow, spr.activeColumn)
		spr.table.SetSelectable(false, false)

		visibility := "v"
		if spr.activeRow > 0 {
			visibility = spr.visibility[spr.keys[spr.activeRow-1]]
		}
		spr.MakeForm(app, visibility)
		modal := CompoundModal(spr.form, 50, 19)
		spr.root = spr.root.AddPage(ModalName, modal, true, true)
		app.SetRoot(spr.root, true)
		app.SetFocus(modal)

	})

	spr.topMenu.AddButton("Delete", func() {
		spr.Hide(spr.activeRow, spr.activeColumn)
		spr.table.SetSelectable(false, false)

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
				spr.UpdateRecord(key, nil, "")
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
		spr.Hide(spr.activeRow, spr.activeColumn)
		spr.table.SetSelectable(false, false)

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
			spr.Hide(spr.activeRow, spr.activeColumn)
			spr.table.SetSelectable(false, false)

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

					app.SetRoot(spr.root, true)
					app.SetFocus(modal)

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
		spr.Hide(spr.activeRow, spr.activeColumn)
		spr.table.SetSelectable(false, false)

		needOldPassword := true
		spr.MakeNewPasswordForm(app, " Change page password ", needOldPassword)
	})
	spr.topMenu.AddButton("Exit", func() {
		saveLabel := spr.topMenu.GetButton(spr.saveMenuInx).GetLabel()
		if saveLabel == "Save" { // nothing to save. Just exit
			clipboard.WriteAll("")
			app.Stop()
			return
		}
		modal := spr.MakeNewModal()
		modal.SetText("Page not saved. Exit?")
		modal.AddButtons([]string{"Exit", "Cancel"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Exit" {
				clipboard.WriteAll("")
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
