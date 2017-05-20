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

createcert:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
