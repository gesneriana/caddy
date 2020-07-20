cd filebrowser
if not exist %cd%/webapp (
		md webapp
    ) else (
		echo %cd%/webapp
    )

tasklist | find /i "filebrowser.exe"
if "%errorlevel%"=="1" (
    echo [ %time:~,-3% ] not exist filebrowser.exe
    cd webapp
    ..\filebrowser.exe -d ..\filebrowser.db
    ) else (
        echo [ %time:~,-3% ] exist filebrowser.exe
    )
