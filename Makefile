.PHONY: dump-schema
dump-schema:
	mysqldump isulibrary -d | grep -v "/\*" | grep -v "\-\-"> webapp/sql/0_schema.sql

.PHONY: dump-data
dump-data:
	mysqldump isulibrary -t | grep -v "/\*" | grep -v "\-\-" | grep -v LOCK | grep -v '^s*$$'> webapp/sql/1_data.sql

.PHONY: init-db
init-db:
	sh webapp/sql/init_db.sh

.PHONY: run-go
run-go:
	cd webapp/go && go run *.go

.PHONY: run-bench
run-bench:
	cd bench && go run main.go wire_gen.go >/dev/null

.PHONY: run-generator
run-generator: init-db
	touch bench/repository/init_data.json
	cd bench/repository/generator && go run *.go
	@make dump-data

.PHONY: up
up:
	cd dev && docker compose up -d

.PHONY: down
down:
	cd dev && docker compose down

.PHONY: compose-log
compose-log:
	cd dev && docker compose logs -f backend
