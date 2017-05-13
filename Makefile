build:
	go build

run:
	./optician-api

buildrun:
	go build && ./optician-api

startdb:
	cockroach start --insecure --host=localhost

cleandb:
	rm -rf cockroach-data optician.db

test:
	go test -v . ./core ./core/store ./core/store/bolt ./core/store/sql
