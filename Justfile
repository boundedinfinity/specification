proj_dir	:= justfile_directory()
m           := "updates"

# Bootstrap macOS
#   brew install just

# Bootstrap Fedora
#   sudo dnf install -y just

list:
    @just --list

all:
	@make iso-3166
	@make iso-639
	@make genc
	@make fips
	@make country-flags
	@make country

iso-3166:
	cd {{ proj_dir }}/iso-3166 && go run iso-3166.go {{ proj_dir }}/iso-3166/iso-3166.tsv {{ proj_dir }}

iso-639:
	cd {{ proj_dir }}/iso-639 && go run iso-639.go {{ proj_dir }}/iso-639.yaml {{ proj_dir }}

genc:
	cd {{ proj_dir }}/genc && go run genc.go "{{ proj_dir }}/genc/GENC Standard Ed1.0.xml" {{ proj_dir }}

fips:
	cd {{ proj_dir }}/fips && go run fips.go {{ proj_dir }}/fips/fips-all.txt {{ proj_dir }}

country-flags:
	cd {{ proj_dir }}/country-flags && go run country-flags.go {{ proj_dir }} {{ proj_dir }}/country-flags/1x1 {{ proj_dir }}/country-flags/4x3

country:
	cd {{ proj_dir }}/country && go run country.go {{ proj_dir }}

mime-types:
	cd {{ proj_dir }}/mime-types && go run mime-types.go {{ proj_dir }}

push:
	git add . || true
	git commit -m {{ m }} || true
	git push origin master
