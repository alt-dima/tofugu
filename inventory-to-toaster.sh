#!/bin/bash

inventory_path=$1

post_to_toaster () {
  # $1 - filename

  dimpath=$1
    dimpath=${dimpath#*$inventory_path}
    dimpath=${dimpath%%.json}
    echo $dimpath

  curl --header "Content-Type: application/json" \
  --request POST \
  -d @$1\
  "$toasterurl/api/dimension/$dimpath?workspace=master&source=inventory&readonly=true" -q
}

for file in $(find $inventory_path -name '*.json')
do
   post_to_toaster $file
done
