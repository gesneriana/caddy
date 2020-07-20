cd ./filebrowser
filebrowser.exe -d ./filebrowser.db config init
filebrowser.exe -d ./filebrowser.db config set --address 0.0.0.0
filebrowser.exe -d ./filebrowser.db config set --port 2021
filebrowser.exe -d ./filebrowser.db config set --locale zh-cn
filebrowser.exe -d ./filebrowser.db config set --log ./filebrowser.log
filebrowser.exe -d ./filebrowser.db users add admin filebrowseradmin --perm.admin
filebrowser.exe -d ./filebrowser.db config set --baseurl /filebrowser