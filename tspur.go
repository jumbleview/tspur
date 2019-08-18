package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

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

// Compound contains all content of the tspur
type spur struct {
	// screens primitives
	flex        *tview.Flex
	lstFlx      *tview.Flex
	list        *tview.List
	form        *tview.Form
	table       *tview.Table
	changeState *tview.Table
	// screen underline data
	keys       []string
	records    map[string][]string
	visibility map[string]string
	width      int
	activeRow  int
	passwd     string
	passwd2    string
}

var scr spur

// CribName is name of the file with csv list
var CribName = "cribsheet.csv"

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
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlue
	tview.Styles.ContrastBackgroundColor = tcell.ColorDarkCyan
	app := tview.NewApplication()

	CribName = cmd[0]
	var sdata []string
	_, errFile := os.Stat(CribName)
	tspurView := func() {
		scr.records = make(map[string][]string)
		scr.visibility = make(map[string]string)
		scr.flex.SetDirection(tview.FlexColumn)
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
		scr.MakeTable(app)
		scr.form = tview.NewForm()
		scr.MakeList(app)
		scr.lstFlx = tview.NewFlex().SetDirection(tview.FlexRow)
		scr.MakeChangeState()
		scr.lstFlx.AddItem(scr.changeState, 1, 0, false)
		scr.lstFlx.AddItem(scr.list, 0, 1, true)
		scr.lstFlx.AddItem(scr.form, 0, 2, false)
		scr.flex.AddItem(scr.lstFlx, 0, 1, true)
		scr.flex.AddItem(scr.table, 0, 2, false)
		app.SetFocus(scr.lstFlx)
	}
	scr.flex = tview.NewFlex()
	scr.flex.SetDirection(tview.FlexRow)
	pwdform := tview.NewForm()
	splash := tview.NewTextView()
	pwdTitle := tview.NewTextView()
	fmt.Fprintf(splash, Splash)
	fmt.Fprintf(pwdTitle, " Enter password:")
	pwdform.SetHorizontal(false)
	pwdInput := func(inp string) {
		scr.passwd = inp
	}
	pwdInput2 := func(inp string) {
		scr.passwd2 = inp
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
		pwdform.AddPasswordField("", "", 14, '*', pwdInput)
		pwdform.AddButton("   submit  ", func() {
			data, err := DecryptFile(CribName, scr.passwd)
			if err == nil {
				if len(data) > 0 {
					sdata = strings.Split(string(data), "\n")
				}
				// clear password form
				//pwdform.Clear(true)
				scr.flex.RemoveItem(pwdform)
				scr.flex.RemoveItem(splash)
				scr.flex.RemoveItem(pwdTitle)
				tspurView()
			} else {
				if data != nil {
					//pwdform.SetTitle(" Wrong Password. Repeat")
					pwdTitle.Clear()
					fmt.Fprintf(pwdTitle, " Wrong Password. Repeat")
					pwdform.Clear(false)
					scr.passwd = ""
					pwdform.AddPasswordField("", "", 14, '*', pwdInput)
					app.SetFocus(scr.flex)
				}
			}
		})
	}
	pwdform.SetBorder(false).SetTitle(" Enter Password:").SetTitleAlign(tview.AlignLeft)
	scr.flex.SetBorder(false).SetTitle(" Enter Password:").SetTitleAlign(tview.AlignLeft)
	scr.flex.AddItem(pwdTitle, 1, 1, false)
	scr.flex.AddItem(pwdform, 5, 1, true)
	scr.flex.AddItem(splash, 0, 1, false)
	app.SetFocus(scr.flex)
	if err := app.SetRoot(scr.flex, true).Run(); err != nil {
		panic(err)
	}
}
