package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

// ConsoleWidth is console horizontal dimension
const ConsoleWidth = 85

// ConsoleHeight is console vertical dimension
const ConsoleHeight = 45

// TopMenuProportion and Table proportional size
const TopMenuProportion = 1
const TableProportion = 7

// SpurTheme is color theme matched to tspur design
type SpurTheme struct {
	MainColor                tcell.Color // Font colors of the table
	MainBackgroundColor      tcell.Color // Background color of te table and top menu
	AccentColor              tcell.Color // Color of top menu font and table borders
	TrackingColor            tcell.Color // Color to illuminate cell hit by ENter key
	FormColor                tcell.Color // Form font and borders color
	FormBackgroundColor      tcell.Color // Form background color
	FormInputBackgroundColor tcell.Color // Form input field background
}

// spur contains all content of the tspur
type Spur struct {
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
	cribPath     string
	cribBase     string
	mode         string
	saveMenuInx  int
	arrowBarrier int
	// to collect comments for commit
	commits []string
	// assigned color theme
	SpurTheme
}

// tspur is cheat sheet table.
// Type of information could be any, but mostly user names and passwords
// Each row consists of key and some values

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
	var tspr Spur

	var GoldenBears = SpurTheme{ // To honor Cal Football team
		MainColor:                tcell.ColorWhite,
		MainBackgroundColor:      tcell.ColorDarkBlue,
		AccentColor:              tcell.ColorGold,
		TrackingColor:            tcell.ColorRed,
		FormColor:                tcell.ColorWhite,
		FormBackgroundColor:      tcell.ColorDarkCyan,
		FormInputBackgroundColor: tcell.ColorDarkBlue,
	}

	tspr.SpurTheme = GoldenBears

	app := tview.NewApplication()

	tspr.root = tview.NewPages()
	tspr.cribName = cmd[0]
	_, errFile := os.Stat(tspr.cribName)
	tspr.cribPath = filepath.Dir(tspr.cribName)
	tspr.cribBase = filepath.Base(tspr.cribName)
	SetDimensions(ConsoleWidth, ConsoleHeight)
	tspr.MakeTopMenu(app)
	tspr.MakeBaseTable(app)
	tspr.flex = tview.NewFlex()
	tspr.flex.SetDirection(tview.FlexRow)
	tspr.flex.SetBorder(false)
	tspr.flex.AddItem(tspr.topMenu, 0, TopMenuProportion, true)
	tspr.flex.AddItem(tspr.table, 0, TableProportion, false)
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
		_, errFile = os.Stat(tspr.cribPath)
		if errFile != nil {
			clipboard.WriteAll("")
			panic(errFile)
		} else {
			needOldPassword := false
			tspr.MakeNewPasswordForm(app, "Create new Page", needOldPassword)
		}
	}
	app.EnableMouse(false)
	if err := app.Run(); err != nil {
		clipboard.WriteAll("")
		panic(err)
	}
}
