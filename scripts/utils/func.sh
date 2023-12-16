#!/bin/bash

is_file_dir(){
	local f="$1"
	[ -f "$f" ] && { echo "$f is a regular file."; exit 0; }
	[ -d "$f" ] && { echo "$f is a directory."; exit 0; }
	[ -L "$f" ] && { echo "$f is a symbolic link."; exit 0; }
	[ -x "$f" ] && { echo "$f is an executeble file."; exit 0; }
}
get_script_path(){
  scriptLocation=$(dirname $0)
  echo "$scriptLocation"
}
get_curr_path(){
  currentPath="$PWD"
  echo "$currentPath"
}
