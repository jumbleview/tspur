package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalName is string assigned to the page with Modal dialog
const ModalName = "modal"

// ModeClipEnter means cell text copied into clipboard when Enter pressed
const ModeClipEnter = "copy-on-enter"

// ModeClipSelect means cell text copied when selected
const ModeClipSelect = "copy-on-select"

// ModeVisibleEnter means cell made visual when Enter pressed
const ModeVisibleEnter = "show-on-enter"

// ModeVisibleSelect means cell mode visual when selected
const ModeVisibleSelect = "show-on-select"

var modeSet [4]string = [4]string{ModeClipEnter, ModeClipSelect, ModeVisibleEnter, ModeVisibleSelect}

// ArrowDefaultBarrier tells at which index to turn on Arrows key converting
const ArrowDefaultBarrier = -1

// ConsoleWidth is console horizontal dimension
const ConsoleWidth = 109

// ConsoleHeight is console vertical dimension
const ConsoleHeight = 45

// TopMenuProportion and Table proportional size
const TopMenuProportion = 1
const TableProportion = 7

// SpurTheme is color theme matched to tspur design
type SpurTheme struct {
	MainColor                tcell.Color // Font colors of the table
	MainBackgroundColor      tcell.Color // Background color of the table and top menu
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
	isLastEventMouse bool
	keys             []string
	records          map[string][]string
	visibility       map[string]string
	width            int
	activeRow        int
	activeColumn     int
	passwd           string
	cribName         string
	cribPath         string
	cribBase         string
	mode             string
	modeIndex        int
	saveMenuInx      int
	arrowBarrier     int
	// to collect comments for commit
	commits []string
	// assigned color theme
	SpurTheme
	// availble modes
	modeSet *[4]string
}

type ColorValues struct {
	Colors     []tcell.Color
	ColorsList string
	Count      int
}

func (v *ColorValues) String() string {
	return v.ColorsList
}

func (v *ColorValues) Set(s string) error {
	if len(s) == 0 {
		return nil
	}
	v.ColorsList = s
	vals := strings.Split(s, ",")
	for _, val := range vals {
		val = strings.Trim(val, " ")
		color := tcell.GetColor(val)
		if color == tcell.ColorDefault {
			return fmt.Errorf("%s is not valid tcell color", val)
		}
		v.Colors = append(v.Colors, color)
	}
	if len(v.Colors) != v.Count {
		return fmt.Errorf("wrong number of colors: %d vs %d", v.Count, len(v.Colors))
	}
	return nil
}

type ModeValue struct {
	Mode       string
	ColorsList string
	Index      int
}

func (v *ModeValue) String() string {
	return v.Mode
}

func (v *ModeValue) Set(s string) error {
	if len(s) == 0 {
		return nil
	}
	for ix, mode := range modeSet {
		if mode == s {
			v.Mode = s
			v.Index = ix
			return nil
		}
	}
	return fmt.Errorf("mode %s is uknown", s)
}

// tspur is cheat sheet table.
// Type of information could be any, but mostly user names and passwords
// Each row consists of key and some values

func main() {
	greeting := "tsupr.exe [-cm] [-cf] [-ct] [-md] [-ta] path_to_data_file"
	var Usage = func() {
		fmt.Fprintln(os.Stderr, greeting)
	}

	var mainColors ColorValues
	mainColors.Count = 2
	flag.Var(&mainColors, "cm", "Colors main: two comma separated  colors: font & background")

	var formColors ColorValues
	formColors.Count = 3
	flag.Var(&formColors, "cf", "Colors formthree comma separated  colors: font, background, & input backgorund")

	var trackingColor ColorValues
	trackingColor.Count = 1
	flag.Var(&trackingColor, "ct", "Color trace: single color: font")

	var tsprMode ModeValue
	var modes []string
	for _, m := range modeSet {
		modes = append(modes, m)
	}
	possibleModes := "Mode: possible values are: " + strings.Join(modes, ",")

	flag.Var(&tsprMode, "md", possibleModes)

	var alterColumn int
	flag.IntVar(&alterColumn, "ta", 0, "Table altering: n > 0 - column  to insert before n; n < 0 - column to delete at -n")

	flag.Parse()
	cmd := flag.Args()
	if len(cmd) != 1 {
		fmt.Fprintf(os.Stderr, "Number of arguments %d\n", len(cmd))
		Usage()
		os.Exit(1)
	}
	var tspr Spur
	tspr.modeSet = &modeSet
	tspr.mode = tsprMode.Mode
	tspr.modeIndex = tsprMode.Index
	var theme = SpurTheme{ // Default is Monochrome theme with red accent for visited cell
		MainColor:                tcell.ColorWhite,
		MainBackgroundColor:      tcell.ColorBlack,
		TrackingColor:            tcell.ColorRed,
		FormColor:                tcell.ColorWhite,
		FormBackgroundColor:      tcell.ColorGray,
		FormInputBackgroundColor: tcell.ColorBlack,
	}

	if len(mainColors.Colors) == 2 {
		theme.MainColor = mainColors.Colors[0]
		theme.MainBackgroundColor = mainColors.Colors[1]
	}
	theme.AccentColor = theme.MainColor

	if len(formColors.Colors) == 3 {
		theme.FormColor = formColors.Colors[0]
		theme.FormBackgroundColor = formColors.Colors[1]
		theme.FormInputBackgroundColor = formColors.Colors[2]
	}

	if len(trackingColor.Colors) == 1 {
		theme.TrackingColor = trackingColor.Colors[0]
	}

	tspr.SpurTheme = theme

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
	tspr.activeRow = 1
	tspr.activeColumn = 1
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		tspr.isLastEventMouse = false
		if tspr.table.HasFocus() {
			if (event.Key() == tcell.KeyRune) && unicode.IsLetter(event.Rune()) {
				runeAsString := strings.ToUpper(string(event.Rune()))
				for ix, ky := range tspr.keys {
					if strings.ToUpper(ky) >= runeAsString {
						tspr.table.Select(ix+1, 1)
						break
					}
				}
				// to suppress letter of navigating the table apart of above
				return nil
			}
		}
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
		tspr.MakeEnterPasswordForm(app, "Enter Password:", alterColumn)
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
	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		tspr.isLastEventMouse = true
		if tspr.topMenu.HasFocus() && action == tview.MouseLeftDown {
			_, eventY := event.Position()
			_, tableY, _, _ := tspr.table.GetRect()
			if eventY > tableY { // mouse clicked in the table area
				tspr.MoveFocusToTable(app)
			}
		}
		return event, action
	})
	// To cleanup clipboard when user closes Window with (x)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		clipboard.WriteAll("")
		os.Exit(0)
	}()
	// To run the App with mouse support enabled
	app.EnableMouse(true)
	if err := app.Run(); err != nil {
		clipboard.WriteAll("")
		panic(err)
	}
}
