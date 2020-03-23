package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ModalName is string assigned to the page with Modal dialog
const ModalName = "modal"

// ModeClipEnter means cell text copied into clipboard when Enter pressed
const ModeClipEnter = "Clipboard-on-Enter"

// ModeClipSelect means cell text copied when selected
const ModeClipSelect = "Clipboard-on-Select"

// ModeVisibleEnter means cell made visual when Enter pressed
const ModeVisibleEnter = "Visible-on-Enter"

// ModeVisibleSelect means cell mode visual when selected
const ModeVisibleSelect = "Visible-on-Select"

// ArrowDefaultBarrier tells at which index to turn on Arrows key converting
const ArrowDefaultBarrier = -1

// spur contains all content of the tspur
type spur struct {
	// tsprs primitives
	root    *tview.Pages // container of pages used in app
	flex    *tview.Flex  // container for the topMenu and the table
	topMenu *tview.Form  // tp menu for the application
	form    *tview.Form  // form used for input/modification of records
	table   *tview.Table // table with records
	modes   *tview.Table // table to select mode

	// tspr underline data
	keys         []string
	records      map[string][]string
	visibility   map[string]string
	width        int
	activeRow    int
	activeColumn int
	passwd       string
	cribName     string
	mode         string
	saveMenuInx  int
	arrowBarrier int
}

// tspur is cheat sheet table.
// Type of infromation could be any, but I keep there my personal user names and passwords
// Each row consists of key and one or more values

func main() {
	greeting := "tsupr.exe path_to_data_file"
	var Usage = func() {
		fmt.Fprintf(os.Stderr, greeting)
	}
	flag.String("-h", "help", greeting)
	flag.Parse()
	cmd := flag.Args()
	if len(cmd) != 1 {
		fmt.Fprintf(os.Stderr, "Number of arguments %d\n", len(cmd))
		Usage()
		os.Exit(1)
	}
	var tspr spur

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDarkBlue
	tview.Styles.PrimaryTextColor = tcell.ColorYellow
	//tview.Styles.SecondaryTextColor = tcell.ColorWhite
	//tview.Styles.TertiaryTextColor = tcell.ColorWhite
	app := tview.NewApplication()

	tspr.root = tview.NewPages()
	tspr.cribName = cmd[0]
	_, errFile := os.Stat(tspr.cribName)

	tspr.records = make(map[string][]string)
	tspr.visibility = make(map[string]string)
	tspr.MakeTopMenu(app)
	tspr.MakeBaseTable(app)
	tspr.flex = tview.NewFlex()
	tspr.flex.SetDirection(tview.FlexRow)
	//scr.flex.SetDirection(tview.FlexColumn)
	tspr.flex.SetBorder(false)
	tspr.flex.AddItem(tspr.topMenu, 0, TopMenuProportion, true)
	tspr.flex.AddItem(tspr.table, 0, 12, false)
	tspr.root = tspr.root.AddPage("table", tspr.flex, true, true)
	app.SetFocus(tspr.flex)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if tspr.topMenu.HasFocus() || tspr.arrowBarrier > ArrowDefaultBarrier {
			_, buttonI := tspr.form.GetFocusedItemIndex()
			if tspr.topMenu.HasFocus() || buttonI >= tspr.arrowBarrier {
				if event.Key() == tcell.KeyRight {
					return tcell.NewEventKey(tcell.KeyTab, 0x09, 0)
				}
				if event.Key() == tcell.KeyLeft {
					return tcell.NewEventKey(tcell.KeyBacktab, 0x09, 0)
				}
			}
			if event.Key() == tcell.KeyDown {
				return tcell.NewEventKey(tcell.KeyTab, 0x09, 0)
			}
			if event.Key() == tcell.KeyUp {
				return tcell.NewEventKey(tcell.KeyBacktab, 0x09, 0)
			}

		}
		if event.Key() == tcell.KeyCtrlC {
			clipboard.WriteAll(event.Name())
		}
		return event
	})
	if errFile == nil {
		tspr.MakeEnterPasswordForm(app, "Enter Password:")
	} else {
		needOldPassword := false
		tspr.MakeNewPasswordForm(app, "Create new Page", needOldPassword)
	}
	if err := app.Run(); err != nil {
		clipboard.WriteAll("")
		panic(err)
	}
}
