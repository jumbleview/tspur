package main

import (
	"fmt"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CompoundModal creates carrier for the Modal Dialog with any primitive (ussally form)
func CompoundModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

// SetFormColors sets colors defining form theme
func SetFormColors(form *tview.Form, background, field, font tcell.Color) {
	form.SetBackgroundColor(background)
	form.SetButtonBackgroundColor(background)
	form.SetFieldBackgroundColor(field)
	form.SetButtonTextColor(font)
	form.SetFieldTextColor(font)
	form.SetLabelColor(font)
}

// MakeForm makes tspr  Form to to insert/modify table record
func (spr *Spur) MakeForm(app *tview.Application, vsbl string) error {
	spr.form = tview.NewForm()
	SetFormColors(spr.form, spr.FormBackgroundColor, spr.FormInputBackgroundColor, spr.FormColor)
	spr.form.SetBorder(true)
	count := spr.width
	if count < 2 {
		count = 2
	}
	var k string
	var v []string
	if spr.activeRow > 0 {
		k = spr.keys[spr.activeRow-1]
	}
	makeInputFields := func(isEmpty bool) {
		spr.form.AddInputField("Record Name", k, 0, nil, func(inp string) {
			k = inp
		})

		if len(k) > 0 && v == nil {
			v = append(v, spr.records[k]...)
		}
		for i := 0; i <= count; i++ {
			valName := "Field " + strconv.Itoa(i+1)
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
			value := v[i]
			if isEmpty {
				value = ""
			}
			if vsbl == "h" {
				spr.form.AddPasswordField(valName, value, 0, '*', changed)
			} else {
				spr.form.AddInputField(valName, value, 0, accepted, changed)
			}
		}
	}
	makeInputFields(false)
	cancel := func() {
		spr.form.Clear(true)
		spr.root.RemovePage(ModalName)
		app.SetFocus(spr.topMenu)
	}

	spr.form.AddButton("Submit", func() {
		if len(k) > 0 {
			_, ok := spr.records[k]
			if !ok {
				spr.keys = append(spr.keys, k)
			}
			for j := len(v) - 1; j >= 0; j-- {
				if len(v[j]) > 0 {
					break
				}
				v = v[:j]
			}
			keyPlace := spr.UpdateRecord(k, v, vsbl)
			spr.UpdateTable(app)
			spr.table.Select(keyPlace+1, 1)
			spr.topMenu.GetButton(spr.saveMenuInx).SetLabel("Save!")
		}

		spr.form.Clear(true)
		spr.root.RemovePage(ModalName)
		spr.MoveFocusToTable(app)
	})
	spr.form.AddButton("Hide/Reveal", func() {

		if vsbl == "h" {
			vsbl = "v"
		} else {
			vsbl = "h"
		}
		spr.form.Clear(false)
		makeInputFields(false)
	})

	spr.form.AddButton("Clear", func() {
		k = ""
		spr.activeRow = -1
		spr.form.Clear(false)
		makeInputFields(true)
	})

	spr.arrowBarrier = spr.form.GetButtonIndex("Submit")
	spr.form.AddButton("Cancel", cancel)

	spr.form.SetCancelFunc(cancel)

	return nil
}

// MakeColumnForm makes tspr  Form to to insert/delete column in table. So far operations at with very first columnimplemented

func (spr *Spur) MakeColumnForm(app *tview.Application, vsbl string) error {
	spr.form = tview.NewForm()
	currentColumn := spr.activeColumn
	if currentColumn < 2 {
		currentColumn = 2
	}
	currentRow := spr.activeRow
	title := "Column:" + strconv.Itoa(currentColumn-1)
	spr.form.SetTitle(title)
	SetFormColors(spr.form, spr.FormBackgroundColor, spr.FormInputBackgroundColor, spr.FormColor)
	spr.form.SetBorder(true)
	count := spr.width
	if count < 2 {
		count = 2
	}
	var v string

	spr.form.AddInputField("Value", "", 21, nil, func(inp string) {
		v = inp
	})

	var sv []string
	insertColumn := func() {
		sv = append(sv, v)
		for k, record := range spr.records {
			record = append(sv, record...)
			spr.records[k] = record
		}
		spr.UpdateTable(app)
		spr.table.Select(currentRow, currentColumn)
		spr.topMenu.GetButton(spr.saveMenuInx).SetLabel("Save!")
	}
	deleteColumn := func() {
		for k, record := range spr.records {
			record = record[1:]
			spr.records[k] = record
		}
		spr.UpdateTable(app)
		spr.table.Select(currentRow, currentColumn)
		spr.topMenu.GetButton(spr.saveMenuInx).SetLabel("Save!")
	}

	cancel := func() {
		spr.form.Clear(true)
		spr.root.RemovePage(ModalName)
		app.SetFocus(spr.topMenu)
	}
	spr.form.AddButton("Insert befor", func() {
		insertColumn()
		cancel()
	})
	spr.form.AddButton("Delete", func() {
		deleteColumn()
		cancel()
	})

	spr.form.AddButton("Cancel", cancel)

	spr.form.SetCancelFunc(cancel)

	return nil
}

func (spr *Spur) Save() {
	csv := ""
	for _, key := range spr.keys {
		line := spr.visibility[key]
		line += ","
		line += key
		values := spr.records[key]
		for _, value := range values {
			line += ","
			line += value
		}
		line += "\n"
		csv += line
	}
	//if len(csv) > 0 {
	err := EncryptFile(spr.cribName, []byte(csv), spr.passwd)
	if err != nil {
		panic(err.Error())
	}
	//}
}

// MakeModeTable makes modal  table to choose Mode
func (spr *Spur) MakeModeTable(app *tview.Application) error {
	//modesSet := [4]string{ModeClipEnter, ModeClipSelect, ModeVisibleEnter, ModeVisibleSelect}
	spr.modes = tview.NewTable().SetBorders(false)

	spr.modes.SetCell(0, 0, tview.NewTableCell(ModeClipEnter).
		SetTextColor(spr.FormColor).SetAlign(tview.AlignCenter).
		SetBackgroundColor(spr.FormBackgroundColor).SetSelectable(true))

	spr.modes.SetCell(1, 0, tview.NewTableCell(ModeClipSelect).
		SetTextColor(spr.FormColor).SetAlign(tview.AlignCenter).
		SetBackgroundColor(spr.FormBackgroundColor).SetSelectable(true))

	spr.modes.SetCell(2, 0, tview.NewTableCell(ModeVisibleEnter).
		SetTextColor(spr.FormColor).SetAlign(tview.AlignCenter).
		SetBackgroundColor(spr.FormBackgroundColor).SetSelectable(true))
	spr.modes.SetCell(3, 0, tview.NewTableCell(ModeVisibleSelect).
		SetTextColor(spr.FormColor).SetAlign(tview.AlignCenter).
		SetBackgroundColor(spr.FormBackgroundColor).SetSelectable(true))

	spr.modes.SetSelectedFunc(func(row, column int) {
		spr.modes.Clear()
		spr.root.RemovePage(ModalName)
		spr.SelectTable(app)
	})
	spr.modes.SetDoneFunc(func(key tcell.Key) {
		if (key == tcell.KeyEnter) || (key == tcell.KeyEscape) {
			spr.modes.Clear()
			spr.root.RemovePage(ModalName)
			// app.SetFocus(spr.topMenu)
			spr.SelectTable(app)
		}
	})
	spr.modes.SetSelectionChangedFunc(func(row, column int) {
		spr.mode = spr.modeSet[row]
		spr.topMenu.GetButton(1).SetLabel(FirstToUpper(spr.mode))
		if !spr.isLastEventMouse {
			return
		}
		spr.modes.Clear()
		spr.root.RemovePage(ModalName)
		//app.SetFocus(spr.topMenu)
		spr.SelectTable(app)
	})

	spr.modes.SetSelectable(true, true)
	var i int
	for i = range spr.modeSet {
		if spr.mode == spr.modeSet[i] {
			break
		}
	}
	spr.modes.Select(i, 0)
	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(spr.modes, 0, 2, true)
	pwdFlex.SetBackgroundColor(spr.FormBackgroundColor)
	pwdFlex.SetTitle("Mode:")
	pwdFlex.SetBorder(true)
	modal := CompoundModal(pwdFlex, len(ModeClipSelect)+3, 6)
	spr.root = spr.root.AddPage(ModalName, modal, true, true)
	app.SetRoot(spr.root, true)
	app.SetFocus(modal)
	spr.arrowBarrier = ArrowDefaultBarrier
	return nil
}

// MakeNewPasswordForm makes tspr  Form to change page password
func (spr *Spur) MakeNewPasswordForm(app *tview.Application, title string, needOldPassword bool) error {
	spr.form = tview.NewForm()
	SetFormColors(spr.form, spr.FormBackgroundColor, spr.FormInputBackgroundColor, spr.FormColor)
	var oldPasswd, passwd1, passwd2 string
	createInputs := func() {
		if needOldPassword {
			spr.form.AddPasswordField("Old Password:", "", 0, '*', func(s string) {
				oldPasswd = s
			})
		}
		spr.form.AddPasswordField("New Password:", "", 0, '*', func(s string) {
			passwd1 = s
		})
		spr.form.AddPasswordField("New Password:", "", 0, '*', func(s string) {
			passwd2 = s
		})
	}
	pwdSubmit := func() {
		if oldPasswd == spr.passwd && passwd1 == passwd2 {
			spr.passwd = passwd1
			spr.Save()
			spr.form.Clear(true)
			if !needOldPassword { // this is case of creating new page
				spr.UpdateTable(app)
			}
			spr.root.RemovePage(ModalName)
			app.SetFocus(spr.topMenu)
		} else {
			spr.form.Clear(true)
			spr.root.RemovePage(ModalName)
			title := " New passwords do not math. Repeat "
			if oldPasswd != spr.passwd {
				title = " Wrong old password. Repeat "
			}
			spr.MakeNewPasswordForm(app, title, needOldPassword)
		}
	}
	createInputs()
	cancel := func() {
		spr.form.Clear(true)
		spr.root.RemovePage(ModalName)
		app.SetFocus(spr.topMenu)
		spr.arrowBarrier = 0
	}
	spr.form.AddButton("Submit", pwdSubmit)
	spr.form.AddButton("Cancel", cancel)
	spr.form.SetCancelFunc(cancel)

	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(spr.form, 0, 2, true)
	pwdFlex.SetBackgroundColor(spr.FormBackgroundColor)

	pwdFlex.SetTitle(title)
	pwdFlex.SetBorder(true) // In case of true border is on black background
	modal := CompoundModal(pwdFlex, 40, 11)
	spr.root = spr.root.AddPage(ModalName, modal, true, true)
	app.SetRoot(spr.root, true)
	app.SetFocus(modal)
	return nil
}

// MakeEnterPasswordForm makes tspr page with Form to enter page password
func (spr *Spur) MakeEnterPasswordForm(app *tview.Application, title string, alterColumn int) error {
	spr.form = tview.NewForm()
	SetFormColors(spr.form, spr.FormBackgroundColor, spr.FormInputBackgroundColor, spr.FormColor)
	var passwd string
	createInputs := func() {
		spr.form.AddPasswordField("", "", 0, '*', func(s string) {
			passwd = s
		})
	}
	pwdSubmit := func() {
		data, err := DecryptFile(spr.cribName, passwd)
		if err == nil {
			spr.AttachData(data, passwd, alterColumn)
			spr.UpdateTable(app)
			spr.activeRow = 1
			spr.activeColumn = 1
			if spr.mode == ModeVisibleSelect {
				spr.Visualize(spr.activeRow, spr.activeColumn)
			} else if spr.mode == ModeClipSelect {
				spr.ToClipBoard(spr.activeRow, spr.activeColumn)
			}
			spr.arrowBarrier = -1

			spr.form.Clear(true)
			spr.root.RemovePage(ModalName)
			spr.MoveFocusToTable(app)
			// app.SetFocus(spr.topMenu)
		} else {
			spr.form.Clear(true)
			spr.root.RemovePage(ModalName)
			title = " Wrong password. Repeat "
			spr.MakeEnterPasswordForm(app, title, alterColumn)
		}
	}
	createInputs()
	cancel := func() {
		app.Stop()
	}
	spr.form.AddButton("Submit", pwdSubmit)
	spr.form.AddButton("Cancel", cancel)
	spr.form.SetCancelFunc(cancel)

	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(spr.form, 0, 2, true)
	pwdFlex.SetBackgroundColor(spr.FormBackgroundColor)
	title = fmt.Sprintf("%s %d", title, alterColumn)
	pwdFlex.SetTitle(title)
	pwdFlex.SetBorder(true)
	modal := CompoundModal(pwdFlex, 27, 7)
	spr.root = spr.root.AddPage(ModalName, modal, true, true)
	app.SetRoot(spr.root, true)
	app.SetFocus(modal)
	return nil
}

// MakeNewModal creates new modal objects and provide its color setting
func (spr *Spur) MakeNewModal() *tview.Modal {
	modal := tview.NewModal()
	modal.SetBackgroundColor(spr.FormBackgroundColor)
	modal.SetButtonBackgroundColor(spr.FormBackgroundColor)
	modal.SetTextColor(spr.FormColor)
	modal.SetButtonTextColor(spr.FormColor)
	return modal
}
