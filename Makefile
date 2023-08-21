migrate:
	@cd migrations && goose mysql "root:1234@/finance" up

migrate-down:
	@cd migrations && goose mysql "root:1234@/finance" down

build:
	@go build

setup:
	@echo "Creating mysql container"
	@docker run --name mysql-container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=1234 -d mysql:latest
	@echo "Waiting for mysql to start"
	@sleep 15
	@echo "Creating database"
	@docker exec -it mysql-container mysql -uroot -p1234 -e "CREATE DATABASE finance"

mysqlsh:
	@docker exec -it mysql-container mysql -uroot -p1234 finance

start-mysql:
	@sudo systemctl stop mysql
	@docker start mysql-container