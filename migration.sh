#!/bin/bash
if [ $1 = "" ]; then
  echo "Migration file name required"
  exit 1
fi
name=$1
timestamp=$(date +"%Y%m%d%H%M")
file_name=$timestamp"_"$name
touch migrations/$file_name.up.sql
touch migrations/$file_name.down.sql