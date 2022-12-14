DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres --network bank_network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
 
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

auth_server:
	go run cmd/auth_server/main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/jiny0x01/simplebank/db/sqlc Store

proto_compile:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb \
		--grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=doc/swagger  --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
		proto/*.proto 

evans:
	evans --port 9090 --host localhost -r repl
	
.PONNY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 db_docs db_schema proto_compile evans auth_server
