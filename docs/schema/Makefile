.PHONY: ojson xlsx

# Convert *.xlsx to *.ojson
ojson:
	./oct convert librarios.xlsx librarios.ojson

# Convert *.ojson to *.xlsx
xlsx:
	./oct convert --notNull=true librarios.ojson librarios.xlsx

# generate *.go file
gorm:
	./oct generate --targetFormat=gorm --package=model --uniqueNameSuffix=_uq --gormModel=BaseModel librarios.ojson output/model.go