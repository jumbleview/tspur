# tspur

Terminal Screen with Protected User Records (TSPUR) is the utility which may serve as terminal cheat sheet or password manager.  

![TSPUR](./images/tspur.png)

It is TUI application. The screen consists of two areas:
* Top menu
* Table with User Records

To start the application supply it with single argument: path to file with data storage. Storage is encoded and password protected. If such a storage does not exists yet application will ask you to enter the password and  will create a new storage.

Application starts with focus on the top menu. Hitting "Enter" while button "Select" is in focus put the focus on the table. To navigate through the table use arrow keys. To put focus  back on top menu use "Esc".

Number of table rows and columns is unlimited but it is unlikely somebody will use more then hundred rows or 3..5 columns. Rows contains one key cell (Record Name), which is always visible and several values which maybe visible or hidden.

Application supports four modes:

![TSPUR_MODE](./images/tspur_mode.png)

* Clipboard-on-Enter. If user hit "Enter" on selected cell  its content is copied into clipboard.

* Clipboard-on-Select. During navigation content of selected cell is copied into clipboard.

* Visible-on-Enter. If user hits "Enter" on selected cell with hidden content it becomes visible.  When cell becomes unselected its contents becomes hidden again. 

* Visible-on-Select. If user selects cell with hidden content it becomes visible. When cell becomes unselected its contents becomes hidden again. 

![TSPUR_SELECT](./images/tspur_select.png)

User may add new cells or edit existing. Row maybe extended with one more value. If there is need to add several values process may be repeated several times.

![TSPUR_EDIT](./images/tspur_edit.png)

# Dependency

Application written in Go language and compiled to standalone executable. All the heavy lifting is done by three imported packages:

	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"github.com/atotto/clipboard"

(Code of this application is just tweaking around "tview" widgets.)

On Linux for clipboard operations to work there is need to install "xclip" .

# Platform Support

Code was developed on Windows 10. Linux (Ubuntu 18.04) seems to be OK as well (just don't forget to install "xclip"). 
For MAC code compiles but I never tried it.

# Mouse Support

Package "tview" does not support mouse and the same is true about this application. 

# Known problem

If size of the terminal windows is changed dynamically structure of the table becomes broken. You can fix it by going back and force between top menu and the table.


