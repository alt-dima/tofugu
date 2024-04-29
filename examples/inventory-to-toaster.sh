#!/bin/bash

inventory_path=$1

post_to_toaster () {
  # $1 - filename

  dimpath=$1
    dimpath=${dimpath#*$inventory_path}
    dimpath=${dimpath%%.json}
    printf "\n-------------- $dimpath -------------\n"

  curl --header "Content-Type: application/json" \
  --request POST \
  -d @$1\
  "$toasterurl/api/dimension/$dimpath?workspace=master&source=inventory&readonly=true" -q
  printf "\n------------------------\n\n"
}

for file in $(find $inventory_path -name '*.json')
do
   post_to_toaster $file
done
