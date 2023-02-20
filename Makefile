.PHONY: dump-schema
dump-schema:
	mysqldump isulibrary -d | grep -v "/\*" | grep -v "\-\-"> webapp/sql/0_schema.sql

.PHONY: dump-data
dump-data:
	mysqldump isulibrary -t | grep -v "/\*" | grep -v "\-\-" | grep -v LOCK> webapp/sql/1_data.sql

.PHONY: init-db
init-db:
	sh webapp/sql/init_db.sh

.PHONY: run-go
run-go:
	cd webapp/go && go run *.go

.PHONY: run-bench
run-bench:
	cd bench && go run *.go >/dev/null

.PHONY: run-generator
run-generator: init-db
	touch bench/repository/init_data.json
	cd bench/repository/generator && go run *.go
	@make dump-data
