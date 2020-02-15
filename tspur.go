package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	//tview.Styles.ContrastBackgroundColor = tcell.ColorDarkCyan
	//c := make(chan os.Signal, 2)
	//	signal.Notify(c, os.Interrupt, syscall.SIGHUP,
	//		syscall.SIGINT,
	//		syscall.SIGTERM,
	//		syscall.SIGQUIT)
	//signal.Notify(c, syscall.SIGTERM)
	//go func() {
	//	select {
	//	case <-c:
	//	clipboard.WriteAll("SIGTERM")
	//	}
	//}()

	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			clipboard.WriteAll(event.Name())
		}
		return event
	})
	scr.root = tview.NewPages()
	scr.cribName = cmd[0]
	var sdata []string
	_, errFile := os.Stat(scr.cribName)
	// tspurView creates main screen with top menu and table
	tspurView := func() {
		scr.records = make(map[string][]string)
		scr.visibility = make(map[string]string)
		scr.flex = tview.NewFlex()
		scr.flex.SetDirection(tview.FlexRow)
		//scr.flex.SetDirection(tview.FlexColumn)
		scr.flex.SetBorder(false)
		for _, s := range sdata {
			// parse string as csv
			elems := strings.Split(s, ",")
			if len(elems) > 1 {
				values := elems[2:]
				scr.keys = append(scr.keys, elems[1])
				if len(values) > scr.width {
					scr.width = len(values)
				}
				scr.records[elems[1]] = values
				scr.visibility[elems[1]] = elems[0]
			}
		}
		scr.MakeTopMenu(app)
		scr.MakeTable(app)
		//scr.root.RemovePage("password")
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

		//app.SetRoot(scr.root, true)
		//scr.form = tview.NewForm()

		//scr.MakeList(app)
		//scr.lstFlx = tview.NewFlex().SetDirection(tview.FlexRow)
		//scr.MakeChangeState()
		//scr.lstFlx.AddItem(scr.changeState, 1, 0, false)
		//scr.lstFlx.AddItem(scr.list, 0, 1, true)
		//scr.lstFlx.AddItem(scr.form, 0, 2, false)
		//scr.flex.AddItem(scr.lstFlx, 0, 1, true)
		//scr.flex.AddItem(scr.table, 0, 2, false)
		//app.SetFocus(scr.lstFlx)
	}

	pwdform := tview.NewForm()
	//splash := tview.NewTextView()
	pwdTitle := tview.NewTextView()
	//fmt.Fprintf(splash, Splash)
	fmt.Fprintf(pwdTitle, " Enter password")
	pwdform.SetHorizontal(false)
	pwdInput := func(inp string) {
		scr.passwd = inp
	}
	pwdInput2 := func(inp string) {
		scr.passwd2 = inp
	}
	// MakePasswordForm creates modal dialog to request a password
	MakePasswordForm := func(prompt string) {
		//scr.root.RemovePage("password")
		pwdTitle := tview.NewTextView()
		fmt.Fprintf(pwdTitle, prompt)
		pwdFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		pwdFlex.AddItem(pwdTitle, 0, 1, false).AddItem(pwdform, 0, 2, true)
		pwdFlex.SetBackgroundColor(tcell.ColorBlue)
		pwdFlex.SetBorder(true) // In case of true border is on black background
		modal := CompoundModal(pwdFlex, 20, 9)
		scr.root = scr.root.AddPage("password", modal, true, true)
		app.SetRoot(scr.root, true)
	}

	if errFile != nil { // path does not exist: create new password
		pwdform.AddPasswordField("1:", "", 14, '*', pwdInput)
		pwdform.AddPasswordField("2:", "", 14, '*', pwdInput2)
		pwdform.AddButton("submit", func() {
			if scr.passwd == scr.passwd2 {
				scr.flex.RemoveItem(pwdform)
				tspurView()
			} else {
				pwdform.SetTitle(" Different passwords. Repeat")
				pwdform.Clear(false)
				pwdform.AddPasswordField("password", "", 14, '*', pwdInput)
				pwdform.AddPasswordField("password", "", 14, '*', pwdInput2)
				app.SetFocus(scr.flex)
			}
		})
	} else { // path does exist: ask old password
		submitLabel := (" submit ")
		pwdform.AddPasswordField(">>", "", len(submitLabel)+2, '*', pwdInput)
		pwdform.AddButton(submitLabel, func() {
			data, err := DecryptFile(scr.cribName, scr.passwd)
			if err == nil {
				if len(data) > 0 {
					sdata = strings.Split(string(data), "\n")
				}
				// clear password form
				//pwdform.Clear(true)
				//scr.flex.RemoveItem(pwdform)
				//scr.flex.RemoveItem(splash)
				//scr.flex.RemoveItem(pwdTitle)
				scr.root.RemovePage("password")
				tspurView()
			} else {
				if data != nil {

					//pwdform.SetTitle(" Wrong Password. Repeat")
					// pwdTitle.Clear()
					// fmt.Fprintf(pwdTitle, " Wrong Password. Repeat")
					scr.root.RemovePage("password")
					pwdform.Clear(false)
					// scr.passwd = ""
					pwdform.AddPasswordField(":", "", len(submitLabel)+2, '*', pwdInput)
					MakePasswordForm(" Wrong Password. Repeat")
					// app.SetFocus(scr.flex)
				}
			}
		})
	}
	MakePasswordForm("Enter Password")
	//app.SetFocus(scr.flex)
	//if err := app.SetRoot(scr.flex, true).Run(); err != nil {
	if err := app.Run(); err != nil {
		clipboard.WriteAll("")
		panic(err)
	}
}
