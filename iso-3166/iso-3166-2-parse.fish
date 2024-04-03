#!/usr/bin/env fish

set script_dir (cd (dirname (status -f)); and pwd)

set name $argv[1]
set url "https://en.wikipedia.org/wiki/$name"
set filename "$script_dir/$name.yaml"
set code_col $argv[2]
set name_col $argv[3]
set category_col $argv[4]

echo "Processing: $url"

set hq_cmd "'{subdivisions: table.sortable > tbody > tr | [{ code: td:nth-child($code_col) > span, name: td:nth-child($name_col), category: td:nth-child($category_col)} ]}'"
set jq_cmd "'del(.subdivisions[] | select(.code == null))'"
set yq_cmd "'.subdivisions'"

echo "curl --silent $url | hq $hq_cmd | jq $jq_cmd | yq -y"
eval "curl --silent $url | hq $hq_cmd | jq $jq_cmd | yq -y" > $filename


