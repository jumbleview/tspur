package main

import (
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// func (scr *spur) MakePasswordForm(app *tview.Application) error {
// 	scr.form = tview.NewForm()
// 	scr.form.AddPasswordField("password", "", 12, '^', func(inp string) {
// 		scr.passwd = inp
// 	})
// 	return nil
// }

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

// MakeForm makes screen  Form to to insert/modify table record
func (scr *spur) MakeForm(app *tview.Application, vsbl string) error {
	scr.form = tview.NewForm()
	scr.form.SetFieldBackgroundColor(tcell.ColorDarkCyan)
	scr.form.SetBorder(true)
	count := scr.width
	var k string
	var v []string
	if len(scr.keys) > 0 {
		if scr.activeRow > 0 {
			k = scr.keys[scr.activeRow-1]
		}
		scr.form.AddInputField("Record Name", k, 21, nil, func(inp string) {
			k = inp
		})
		if len(k) > 0 {
			v = append(v, scr.records[k]...)
		}
		// if vsbl == "" {
		// 	vsbl = scr.visibility[k]
		// }
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
	} else {
		v = append(v, "")
		scr.form.AddInputField("Record Name", "", 21, nil, func(inp string) {
			k = inp
		})
		scr.form.AddInputField("+", "", 21, nil, func(inp string) {
			v[0] = inp
		})
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
			//scr.records[k] = v
			//scr.visibility[k] = presentation
			keyPlace := scr.UpdateRecords(k, v, presentation)
			scr.UpdateTable(app)
			scr.table.Select(keyPlace+1, 1)
			scr.topMenu.GetButton(4).SetLabel("Save!")
			//scr.ChangeState(StateAlert)
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
		scr.table.SetSelectable(true, true)
		app.SetFocus(scr.table)
	})
	scr.form.AddButton("Save visible", func() {
		submit("v")
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		scr.table.SetSelectable(true, true)
		app.SetFocus(scr.table)
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
	if len(csv) > 0 {
		err := EncryptFile(scr.cribName, []byte(csv), scr.passwd)
		if err != nil {
			panic(err.Error())
		}
		//scr.ChangeState(StateSaved)
	}
	//app.Stop()
	//lnavigate()
}

// MakeSaveForm makes screen  Form to apporve saving of the changed page
func (scr *spur) MakeSaveForm(app *tview.Application, vsbl string) error {
	scr.form = tview.NewForm()
	var dropDown []string
	dropDown = append(dropDown, scr.cribName+" ?")
	scr.form.AddDropDown("Save the page:", dropDown, 0, nil)
	scr.form.AddButton("Save", func() {
		scr.Save()
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		app.SetFocus(scr.topMenu)
	})
	scr.form.AddButton("Cancel", func() {
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		app.SetFocus(scr.topMenu)
	})
	return nil
}

// MakeSaveForm makes screen  Form to apporve saving of the changed page
func (scr *spur) MakeNewPasswordForm(app *tview.Application) error {
	scr.form = tview.NewForm()
	scr.form.AddPasswordField("Old Password:", "", 21, '*', nil)
	scr.form.AddPasswordField("New Password:", "", 21, '*', nil)
	scr.form.AddPasswordField("New Password:", "", 21, '*', nil)

	scr.form.AddButton("Submit", func() {
		scr.Save()
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		app.SetFocus(scr.topMenu)
	})
	scr.form.AddButton("Cancel", func() {
		scr.form.Clear(true)
		scr.root.RemovePage(ModalName)
		app.SetFocus(scr.topMenu)
	})
	return nil
}
