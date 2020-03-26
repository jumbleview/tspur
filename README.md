# tspur

Terminal Screen with Protected User Records (TSPUR) is the utility which put under the tip of fingers important personal information. It may serve as password manager or as terminal cheat sheet.  

![TSPUR](./images/tspur.png)

The screen consists of two areas:
* Top Menu
* Table with User Records
To start the application supply it with single argument: path to file with data storage. Storage is encoded and password protected. If such a storage does not exists yet application asks you to supply the password and amy create new storage for you

Application starts with focus on Top Menu. Hitting "Enter" while button "Select" in focus switches focus to the Table.
To navigate through the table use arrow keys. To switch back to top menu hit "Esc".

Number of table rows and columns is unlimited but it is unlikely somebody with use more then hundred rows or 3..5 columns. 
Rows contains one key cell, which is always visible and several values which maybe visible or hidden.


![TSPUR_MODE](./images/tspur_mode.png)

Application supports four modes:
* Clipboard-on-Enter. When user hit "Enter" while some cell is selected its content is copied into clipboard.

* Clipboard-on-Select. When user navigate to some cell its content is copied into clipboard.

* Visible-on-Enter. When user hits "Enter" while some cell is selected its hidden content becomes visible. There is no effect on visible cell. When cell becomes unselected its contents becomes hidden again. 

* Visible-on-Select. When user selects cell with hidden content it becomes visible. When cell becomes unselected its contents becomes hidden again. 

![TSPUR_SELECT](./images/tspur_select.png)

User may add new cell or edit existing. Row maybe extended with one more value. If there is need to add several values process must be repeated several times.

![TSPUR_EDIT](./images/tspur_edit.png)





