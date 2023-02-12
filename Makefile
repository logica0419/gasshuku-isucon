.PHONY: dump-schema
dump-schema:
	mysqldump isulibrary -d | grep -v "/\*" | grep -v "\-\-"> webapp/sql/0_schema.sql

.PHONY: dump-data
dump-data:
	mysqldump isulibrary -t | grep -v "/\*" | grep -v "\-\-" | grep -v LOCK> webapp/sql/1_data.sql

.PHONY: run-go
run-go:
	cd webapp/go && go run *.go
