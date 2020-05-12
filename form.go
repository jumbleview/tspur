package main

import (
	"strconv"

	// "github.com/atotto/clipboard"
	"github.com/d-tsuji/clipboard"
	"github.com/gdamore/tcell"
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
	spr.form.AddInputField("Record Name", k, 21, nil, func(inp string) {
		k = inp
	})
	if len(k) > 0 {
		v = append(v, spr.records[k]...)
	}
	for i := 0; i <= count; i++ {
		valName := "Field " + strconv.Itoa(i)
		if i == count {
			valName = "+"
		}
		if i >= len(v) {
			v = append(v, "")
		}
		locali := i
		accepted := func(inp string, last rune) bool {
			//clipboard.WriteAll(inp)
			clipboard.Set(inp)
			return true
		}
		changed := func(inp string) {
			v[locali] = inp
			//clipboard.WriteAll(inp)
			clipboard.Set(inp)
		}
		if vsbl == "h" {
			spr.form.AddPasswordField(valName, v[i], 21, '*', changed)
		} else {
			spr.form.AddInputField(valName, v[i], 21, accepted, changed)
		}
	}
	submit := func(presentation string) {
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
			keyPlace := spr.UpdateRecords(k, v, presentation)
			spr.UpdateTable(app)
			spr.table.Select(keyPlace+1, 1)
			spr.topMenu.GetButton(spr.saveMenuInx).SetLabel("Save!")
		}
	}
	cancel := func() {
		spr.form.Clear(true)
		spr.root.RemovePage(ModalName)
		app.SetFocus(spr.topMenu)
	}
	spr.form.AddButton("Save hidden", func() {
		submit("h")
		spr.form.Clear(true)
		spr.root.RemovePage(ModalName)
		spr.MoveFocusToTable(app)
	})
	spr.arrowBarrier = spr.form.GetButtonIndex("Save hidden")
	spr.form.AddButton("Save visible", func() {
		submit("v")
		spr.form.Clear(true)
		spr.root.RemovePage(ModalName)
		spr.MoveFocusToTable(app)
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
	modesSet := [4]string{ModeClipEnter, ModeClipSelect, ModeVisibleEnter, ModeVisibleSelect}
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
		spr.mode = modesSet[row]
		spr.modes.Clear()
		spr.root.RemovePage(ModalName)
		app.SetFocus(spr.topMenu)
	})
	spr.modes.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
		} else if key == tcell.KeyEscape {
			spr.modes.Clear()
			spr.root.RemovePage(ModalName)
			app.SetFocus(spr.topMenu)
		}
	})
	spr.modes.SetSelectable(true, true)
	var i int
	for i = range modesSet {
		if spr.mode == modesSet[i] {
			break
		}
	}
	spr.modes.Select(i, 0)
	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(spr.modes, 0, 2, true)
	pwdFlex.SetBackgroundColor(spr.FormBackgroundColor)
	pwdFlex.SetTitle("Mode:")
	pwdFlex.SetBorder(true)
	modal := CompoundModal(pwdFlex, 21, 6)
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
			spr.form.AddPasswordField("Old Password:", "", 21, '*', func(s string) {
				oldPasswd = s
			})
		}
		spr.form.AddPasswordField("New Password:", "", 21, '*', func(s string) {
			passwd1 = s
		})
		spr.form.AddPasswordField("New Password:", "", 21, '*', func(s string) {
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
func (spr *Spur) MakeEnterPasswordForm(app *tview.Application, title string) error {
	spr.form = tview.NewForm()
	SetFormColors(spr.form, spr.FormBackgroundColor, spr.FormInputBackgroundColor, spr.FormColor)
	var passwd string
	createInputs := func() {
		spr.form.AddPasswordField("", "", 21, '*', func(s string) {
			passwd = s
		})
	}
	pwdSubmit := func() {
		data, err := DecryptFile(spr.cribName, passwd)
		if err == nil {
			spr.AttachData(data, passwd)
			spr.UpdateTable(app)
			spr.form.Clear(true)
			spr.root.RemovePage(ModalName)
			app.SetFocus(spr.topMenu)
		} else {
			spr.form.Clear(true)
			spr.root.RemovePage(ModalName)
			title = " Wrong password. Repeat "
			spr.MakeEnterPasswordForm(app, title)
		}
	}
	createInputs()
	cancel := func() {
		app.Stop()
		return
	}
	spr.form.AddButton("Submit", pwdSubmit)
	spr.form.AddButton("Cancel", cancel)
	spr.form.SetCancelFunc(cancel)

	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(spr.form, 0, 2, true)
	pwdFlex.SetBackgroundColor(spr.FormBackgroundColor)
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
