#!/bin/bash
cd ./filebrowser
./filebrowser -d ./filebrowser.db config init
./filebrowser -d ./filebrowser.db config set --address 0.0.0.0
./filebrowser -d ./filebrowser.db config set --port 2021
./filebrowser -d ./filebrowser.db config set --locale zh-cn
./filebrowser -d ./filebrowser.db config set --log ./filebrowser.log
./filebrowser -d ./filebrowser.db users add admin filebrowseradmin --perm.admin
./filebrowser -d ./filebrowser.db config set --baseurl /filebrowser
