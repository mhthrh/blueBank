dbConnection= postgres://postgres:123456@localhost:5432/blue_bank?sslmode=disable
outPutPath=C:\Users\Mohsen\Desktop\RunApp

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
	migrate -path Db/Migration -database "$(dbConnection)" -verbose up 2

migrate_down:
	migrate -path Db/Migration -database "$(dbConnection)" -verbose down 1

migrate_command:
	migrate create -ext sql -dir Db/Migration/ -seq initial

create_pb:
	protoc --go-grpc_out=Proto  --go_out=Proto  Proto/*.proto

encrypt_config:
	  go run ./ConfigFiles/main.go

go_test:
	go test -v -cover ./Test/...

build:
	cd  "$(outPutPath)"
	del "$(outPutPath)" *.exe
	go build -o "$(outPutPath)"/GrpcServices.exe ./GrpcApi/main.go
	go build -o "$(outPutPath)"/ApiGateway.exe ./ApiGateway/main.go
	go build -o "$(outPutPath)"/Dispatcher.exe ./Dispatcher/main.go
