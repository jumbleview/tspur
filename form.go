package main

import (
	"strconv"

	"github.com/atotto/clipboard"
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
func (scr *spur) MakeForm(app *tview.Application, vsbl string) error {
	scr.form = tview.NewForm()
	SetFormColors(scr.form, tcell.ColorDarkCyan, tcell.ColorDarkBlue, tcell.ColorWhite)
	scr.form.SetBorder(true)
	count := scr.width
	if count < 2 {
		count = 2
	}
	var k string
	var v []string
	if scr.activeRow > 0 {
		k = scr.keys[scr.activeRow-1]
	}
	scr.form.AddInputField("Record Name", k, 21, nil, func(inp string) {
		k = inp
	})
	if len(k) > 0 {
		v = append(v, scr.records[k]...)
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
			clipboard.WriteAll(inp)
			return true
		}
		changed := func(inp string) {
			v[locali] = inp
			clipboard.WriteAll(inp)
		}
		if vsbl == "h" {
			scr.form.AddPasswordField(valName, v[i], 21, '*', changed)
		} else {
			scr.form.AddInputField(valName, v[i], 21, accepted, changed)
		}
	}
	submit := func(presentation string) {
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
			keyPlace := scr.UpdateRecords(k, v, presentation)
			scr.UpdateTable(app)
			scr.table.Select(keyPlace+1, 1)
			scr.topMenu.GetButton(scr.saveMenuInx).SetLabel("Save!")
		}
	}
	cancel := func() {
		scr.form.Clear(true)
		//cell := scr.table.GetCell(scr.table.GetSelection())
		//clipboard.WriteAll(cell.Text)
		scr.root.RemovePage(ModalName)
		app.SetFocus(scr.topMenu)
	}
	scr.form.AddButton("Save hidden", func() {
		submit("h")
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		scr.MoveFocusToTable(app)
	})
	scr.arrowBarrier = scr.form.GetButtonIndex("Save hidden")
	scr.form.AddButton("Save visible", func() {
		submit("v")
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		scr.MoveFocusToTable(app)
	})

	scr.form.AddButton("Cancel", cancel)

	scr.form.SetCancelFunc(cancel)

	return nil
}

func (scr *spur) Save() {
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
	//if len(csv) > 0 {
	err := EncryptFile(scr.cribName, []byte(csv), scr.passwd)
	if err != nil {
		panic(err.Error())
	}
	//}
}

// MakeModeTable makes modal  table to choose Mode
func (scr *spur) MakeModeTable(app *tview.Application) error {
	modesSet := [4]string{ModeClipEnter, ModeClipSelect, ModeVisibleEnter, ModeVisibleSelect}
	scr.modes = tview.NewTable().SetBorders(false)

	scr.modes.SetCell(0, 0, tview.NewTableCell(ModeClipEnter).
		SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDarkCyan).SetSelectable(true))

	scr.modes.SetCell(1, 0, tview.NewTableCell(ModeClipSelect).
		SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDarkCyan).SetSelectable(true))

	scr.modes.SetCell(2, 0, tview.NewTableCell(ModeVisibleEnter).
		SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDarkCyan).SetSelectable(true))
	scr.modes.SetCell(3, 0, tview.NewTableCell(ModeVisibleSelect).
		SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorDarkCyan).SetSelectable(true))

	scr.modes.SetSelectedFunc(func(row, column int) {
		scr.mode = modesSet[row]
		scr.modes.Clear()
		scr.root.RemovePage(ModalName)
		app.SetFocus(scr.topMenu)
	})
	scr.modes.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
		} else if key == tcell.KeyEscape {
			scr.modes.Clear()
			scr.root.RemovePage(ModalName)
			app.SetFocus(scr.topMenu)
		}
	})
	scr.modes.SetSelectable(true, true)
	var i int
	for i = range modesSet {
		if scr.mode == modesSet[i] {
			break
		}
	}
	scr.modes.Select(i, 0)
	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(scr.modes, 0, 2, true)
	pwdFlex.SetBackgroundColor(tcell.ColorDarkCyan)
	pwdFlex.SetTitle("Mode:")
	pwdFlex.SetBorder(true)
	modal := CompoundModal(pwdFlex, 21, 6)
	scr.root = scr.root.AddPage(ModalName, modal, true, true)
	app.SetRoot(scr.root, true)
	app.SetFocus(modal)
	scr.arrowBarrier = ArrowDefaultBarrier
	return nil
}

// MakeNewPasswordForm makes tspr  Form to change page password
func (scr *spur) MakeNewPasswordForm(app *tview.Application, title string, needOldPassword bool) error {
	scr.form = tview.NewForm()
	SetFormColors(scr.form, tcell.ColorDarkCyan, tcell.ColorDarkBlue, tcell.ColorWhite)
	var oldPasswd, passwd1, passwd2 string
	createInputs := func() {
		if needOldPassword {
			scr.form.AddPasswordField("Old Password:", "", 21, '*', func(s string) {
				oldPasswd = s
			})
		}
		scr.form.AddPasswordField("New Password:", "", 21, '*', func(s string) {
			passwd1 = s
		})
		scr.form.AddPasswordField("New Password:", "", 21, '*', func(s string) {
			passwd2 = s
		})
	}
	pwdSubmit := func() {
		if oldPasswd == scr.passwd && passwd1 == passwd2 {
			scr.passwd = passwd1
			scr.Save()
			scr.form.Clear(true)
			if !needOldPassword { // this is case of creating new page
				scr.UpdateTable(app)
			}
			scr.root.RemovePage(ModalName)
			app.SetFocus(scr.topMenu)
		} else {
			scr.form.Clear(true)
			scr.root.RemovePage(ModalName)
			title := " New passwords do not math. Repeat "
			if oldPasswd != scr.passwd {
				title = " Wrong old password. Repeat "
			}
			scr.MakeNewPasswordForm(app, title, needOldPassword)
		}
	}
	createInputs()
	cancel := func() {
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		app.SetFocus(scr.topMenu)
		scr.arrowBarrier = 0
	}
	scr.form.AddButton("Submit", pwdSubmit)
	scr.form.AddButton("Cancel", cancel)
	scr.form.SetCancelFunc(cancel)

	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(scr.form, 0, 2, true)
	pwdFlex.SetBackgroundColor(tcell.ColorDarkCyan)
	pwdFlex.SetTitle(title)
	pwdFlex.SetBorder(true) // In case of true border is on black background
	modal := CompoundModal(pwdFlex, 40, 11)
	scr.root = scr.root.AddPage(ModalName, modal, true, true)
	app.SetRoot(scr.root, true)
	app.SetFocus(modal)
	return nil
}

// MakeEnterPasswordForm makes tspr page with Form to enter page password
func (scr *spur) MakeEnterPasswordForm(app *tview.Application, title string) error {
	scr.form = tview.NewForm()
	SetFormColors(scr.form, tcell.ColorDarkCyan, tcell.ColorDarkBlue, tcell.ColorWhite)
	var passwd string
	createInputs := func() {
		scr.form.AddPasswordField("", "", 21, '*', func(s string) {
			passwd = s
		})
	}
	pwdSubmit := func() {
		data, err := DecryptFile(scr.cribName, passwd)
		if err == nil {
			scr.AttachData(data, passwd)
			scr.UpdateTable(app)
			scr.form.Clear(true)
			scr.root.RemovePage(ModalName)
			app.SetFocus(scr.topMenu)
		} else {
			scr.form.Clear(true)
			scr.root.RemovePage(ModalName)
			title = " Wrong password. Repeat "
			scr.MakeEnterPasswordForm(app, title)
		}
	}
	createInputs()
	cancel := func() {
		app.Stop()
		return
	}
	scr.form.AddButton("Submit", pwdSubmit)
	scr.form.AddButton("Cancel", cancel)
	scr.form.SetCancelFunc(cancel)

	pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	pwdFlex.AddItem(scr.form, 0, 2, true)
	pwdFlex.SetBackgroundColor(tcell.ColorDarkCyan)
	pwdFlex.SetTitle(title)
	pwdFlex.SetBorder(true)
	modal := CompoundModal(pwdFlex, 27, 7)
	scr.root = scr.root.AddPage(ModalName, modal, true, true)
	app.SetRoot(scr.root, true)
	app.SetFocus(modal)
	return nil
}
