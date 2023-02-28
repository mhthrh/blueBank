Database_connection= postgres://postgres:123456@localhost:5432/blue_bank?sslmode=disable

compose_up:
	docker-compose up -d

compose_stop:
	docker-compose stop

compose_down:
	docker-compose down

create_db:
	docker exec -it postgres psql -U postgres -c "create database blue_bank"

drop_db:
	docker exec -it postgres psql -U postgres -c "drop database blue_bank with (force)"

migrate_up:
	migrate -path Db/Migration -database "$(Database_connection)" -verbose up 2

migrate_down:
	migrate -path Db/Migration -database "$(Database_connection)" -verbose down 1

migrate_command:
	migrate create -ext sql -dir Db/Migration/ -seq initial

create_pb:
	protoc --go-grpc_out=Proto  --go_out=Proto  Proto/*.proto

encrypt_config:
	  go run ./ConfigFiles/main.go

go_test:
	go test -v -cover ./Test/...
