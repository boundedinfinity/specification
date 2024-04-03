#!/usr/bin/env fish

set script_dir (cd (dirname (status -f)); and pwd)
set url https://en.wikipedia.org/wiki/ISO_3166-1#Codes

set name_col 1
set alpha_2_col 2
set alpha_3_col 3
set num_col 4
set link_col 5
set filename  $script_dir/iso-3166-1.yml

echo "Processing: $url"

set hq_cmd "'{ entries: table.sortable caption + tbody > tr | [{
    aname: td:nth-child($name_col),
    codealpha2: td:nth-child($alpha_2_col),
    codealpha3: td:nth-child($alpha_3_col),
    codenumeric: td:nth-child($num_col),
    link: td:nth-child($link_col)
}] }'"

set jq_cmd "'del(.entries[] | select(.aname == null))'"
set yq_cmd "'.entries'"

set cmd "curl --silent $url | hq $hq_cmd | jq $jq_cmd | yq -y"

echo $cmd
eval $cmd > $filename



sed -i '' 's/aname/name/g' $filename
sed -i '' 's/codealpha2/code-alpha-2/g' $filename
sed -i '' 's/codealpha3/code-alpha-3/g' $filename
sed -i '' 's/codenumeric/code-numeric/g' $filename
sed -i '' -E 's/ ([0-9]{3})$/ \'\1\'/g' $filename
sed -i '' 's/ISO 3166-2/ISO_3166-2/g' $filename
sed -i '' -E 's/\[.*\]//g' $filename

# for link in (cat $filename | yq -r '.entries[].link')
#     echo "https://en.wikipedia.org/wiki/$link"
# end


$script_dir/iso-3166-2-parse.fish ISO_3166-2:AF 1 3 -1
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:AX -1 -1 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AL 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:DZ 1 2 -1
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:AS -1 -1 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AD 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AO 1 2 -1
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:AI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:AQ
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AG 1 2 3
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AR 1 2 3
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AM 1 2 5
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:AW
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AU 1 2 3
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AT 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:AZ 1 2 4
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BS 1 2 3
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BH 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BD 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BB 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BY 1 8 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BE 1 3 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BZ 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BJ 1 2 -1
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:BM
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BT 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BO 1 2 3
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BQ 1 2 -1
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:BA 1 3 4
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BW 1 2 3
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:BV
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BR 1 2 3
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IO
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BN 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BG 1 2 -1
$script_dir/iso-3166-2-parse.fish ISO_3166-2:BF 1 2 -1
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:BI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CV
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CF
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TD
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CX
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CD
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:HR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CW
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CZ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:DK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:DJ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:DM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:DO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:EC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:EG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SV
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GQ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ER
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:EE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SZ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ET
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:FK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:FO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:FJ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:FI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:FR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GF
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PF
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TF
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:DE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GD
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GP
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GW
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:HT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:HM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:VA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:HN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:HK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:HU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ID
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IQ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:IT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:JM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:JP
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:JE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:JO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KZ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KP
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KW
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LV
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LB
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MW
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MV
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ML
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MQ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:YT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MX
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:FM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MD
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ME
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MZ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NP
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NZ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NF
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MP
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:NO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:OM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PW
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:QA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:RE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:RO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:RU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:RW
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:BL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:KN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:MF
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:PM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:VC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:WS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ST
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:RS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SX
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SB
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ZA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SS
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ES
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:LK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SD
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SJ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:CH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:SY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TW
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TJ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TZ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TL
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TK
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TO
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TT
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TR
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TC
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:TV
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:UG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:UA
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:AE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:GB
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:US
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:UM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:UY
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:UZ
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:VU
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:VE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:VN
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:VG
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:VI
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:WF
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:EH
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:YE
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ZM
# $script_dir/iso-3166-2-parse.fish ISO_3166-2:ZW
