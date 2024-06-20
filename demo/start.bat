rem Using conhost is optional, but it allows proper set of terminal size out of programm to match file with data

rem start with default setting
start conhost ..\tspur.exe cribsheet

rem  color theme monochrome dark (the same as default)
rem start conhost ..\tspur.exe -cm=white,black -cf=white,gray,black -ct=red -md=copy-on-enter cribsheet

rem  color theme  monochrome light 
rem start conhost ..\tspur.exe -cm=black,silver -cf=black,silver,white -ct=red -md=copy-on-select cribsheet

rem color theme aka Norton
rem start conhost ..\tspur.exe -cm=white,darkblue -cf=white,darkcyan,darkblue -ct=red -md=copy-on-enter cribsheet

rem color theme green
rem start conhost ..\tspur.exe -cm=silver,darkseagreen -cf=silver,darkseagreen,lavender -ct=red -md=copy-on-enter cribsheet


