makefile_dir	:= $(abspath $(shell pwd))
m				:= "updates"

list:
	@grep '^[^#[:space:]].*:' Makefile | grep -v ':=' | grep -v '^\.' | sed 's/:.*//g' | sed 's/://g' | sort

.PHONY: all
all:
	@make iso-3166
	@make iso-639
	@make genc
	@make fips
	@make country-flags
	@make country

.PHONY: iso-3166
iso-3166:
	cd $(makefile_dir)/iso-3166 && go run iso-3166.go $(makefile_dir)/iso-3166/iso-3166.tsv $(makefile_dir)

.PHONY: iso-639
iso-639:
	cd $(makefile_dir)/iso-639 && go run iso-639.go $(makefile_dir)/iso-639.yaml $(makefile_dir)

.PHONY: genc
genc:
	cd $(makefile_dir)/genc && go run genc.go "$(makefile_dir)/genc/GENC Standard Ed1.0.xml" $(makefile_dir)

.PHONY: fips
fips:
	cd $(makefile_dir)/fips && go run fips.go $(makefile_dir)/fips/fips-all.txt $(makefile_dir)

.PHONY: country-flags
country-flags:
	cd $(makefile_dir)/country-flags && go run country-flags.go $(makefile_dir) $(makefile_dir)/country-flags/1x1 $(makefile_dir)/country-flags/4x3

.PHONY: country
country:
	cd $(makefile_dir)/country && go run country.go $(makefile_dir)
