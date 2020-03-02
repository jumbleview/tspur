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

// StateSaved is what to show when new data is in Table
const StateSaved = "== Saved =="

// StateAlert is waht to show when table is saved
const StateAlert = "~ Changed ~"

// Splash is visiting screen greeting
const Splash = "" +
	"                               SSS\n" +
	"                              S   S\n" +
	"                             S     S\n" +
	"                             S\n" +
	"                   t         S\n" +
	"                   t          S\n" +
	"                   t            S\n" +
	"                 ttttt           S         pppppp     u     u    r  r r r\n" +
	"                   t               S       p     p    u     u    r r r r  \n" +
	"                   t                S      p     p    u     u    r        \n" +
	"                   t                S      p     p    u     u    r        \n" +
	"                   t          S     S      p     p    u     u    r        \n" +
	"                   t  t        S   S       p     p    u     u    r        \n" +
	"                    tt          SSS        pppppp      uuuuu u   r        \n" +
	"                                           p                              \n" +
	"                                           p                              \n" +
	"                                           p                              \n" +
	"                                           p                              \n"

// spur contains all content of the tspur
type spur struct {
	// screens primitives
	root    *tview.Pages // container of pages used in app
	flex    *tview.Flex  // container for the topMenu and the table
	topMenu *tview.Form  // tp menu for the application
	form    *tview.Form  // form used for input/modification of records
	table   *tview.Table // table with records
	// to be deleted

	lstFlx      *tview.Flex
	list        *tview.List
	changeState *tview.Table
	// screen underline data
	keys         []string
	records      map[string][]string
	visibility   map[string]string
	width        int
	activeRow    int
	activeColumn int
	passwd       string
	passwd2      string
	cribName     string
}

var scr spur

// tspur is cheat sheet table.
// Type of infromation could be any.
// Each row consists of key and one or more values

func main() {
	// Read cribsheet file and present it as the table
	//data, _ := ioutil.ReadFile(CribName)
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
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDarkBlue
	tview.Styles.PrimaryTextColor = tcell.ColorYellow

	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			clipboard.WriteAll(event.Name())
		}
		return event
	})
	scr.root = tview.NewPages()
	scr.cribName = cmd[0]
	//var sdata []string
	_, errFile := os.Stat(scr.cribName)

	scr.records = make(map[string][]string)
	scr.visibility = make(map[string]string)
	scr.MakeTopMenu(app)
	scr.MakeBaseTable(app)
	scr.flex = tview.NewFlex()
	scr.flex.SetDirection(tview.FlexRow)
	//scr.flex.SetDirection(tview.FlexColumn)
	scr.flex.SetBorder(false)
	scr.flex.AddItem(scr.topMenu, 0, 1, true)
	scr.flex.AddItem(scr.table, 0, 12, false)
	scr.root = scr.root.AddPage("table", scr.flex, true, true)
	app.SetFocus(scr.flex)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if !scr.table.HasFocus() {
			if event.Key() == tcell.KeyRight || event.Key() == tcell.KeyDown {
				return tcell.NewEventKey(tcell.KeyTab, 0x09, 0)
			}
			if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyUp {
				return tcell.NewEventKey(tcell.KeyBacktab, 0x09, 0)
			}
		}
		return event
	})
	if errFile == nil {
		scr.MakeEnterPasswordForm(app, "Enter Password:")
	} else {
		needOldPassword := false
		scr.MakeNewPasswordForm(app, "Create new Page", needOldPassword)
	}
	if err := app.Run(); err != nil {
		clipboard.WriteAll("")
		panic(err)
	}
}
