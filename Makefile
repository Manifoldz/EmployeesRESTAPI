
PASSWORD=qwerty
HOST_PORT=5432

CONTAINER_DB_NAME=employeesrestapi-db-1
CONTAINER_APP_NAME=employeesrestapi-employeessrv-1

all:

build: 
	docker-compose build employeessrv

run:
	docker-compose up employeessrv

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up

docker-run:
	docker run --name=$(CONTAINER_DB_NAME) -e POSTGRES_PASSWORD=$(PASSWORD) -p $(HOST_PORT):5432 -d postgres

docker-rm:
	docker stop $(CONTAINER_DB_NAME)
	docker rm $(CONTAINER_DB_NAME)

migrate-up:
	migrate -path ./schema -database \
	'postgres://postgres:$(PASSWORD)@localhost:$(HOST_PORT)/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database \
	'postgres://postgres:$(PASSWORD)@localhost:$(HOST_PORT)/postgres?sslmode=disable' down

connect-db:
	docker exec -it $(CONTAINER_DB_NAME) /bin/bash

connect-app:
	docker exec -it $(CONTAINER_DB_NAME) /bin/bash

test1: post-createEmployees
test2: get-AllEmployees get-AllEmployees1 get-AllEmployees1Develop1
test3: delete-someEmployees get-AllEmployees
test4: put-updateEmployee get-AllEmployees

# Создать сотрудников
post-createEmployees:
	curl -i -X POST localhost:8000/api/employees/ \
	-H "Content-Type: application/json" \
	-d '{"id":1,"name":"Aleks","surname":"First","phone":"11111111","company_id":1,"passport":{"type":"RF","number":"1111"},"department":{"name":"Develop1","phone":"1111"}}';
	curl -i -X POST localhost:8000/api/employees/ \
	-H "Content-Type: application/json" \
	-d '{"id":2,"name":"Sergei","surname":"Second","phone":"22222222","company_id":1,"passport":{"type":"RF","number":"2222"},"department":{"name":"Develop1","phone":"1111"}}';
	curl -i -X POST localhost:8000/api/employees/ \
	-H "Content-Type: application/json" \
	-d '{"id":3,"name":"Toya","surname":"Third","phone":"333333333","company_id":3,"passport":{"type":"driver","number":"33333"},"department":{"name":"Market","phone":"33333"}}';
	curl -i -X POST localhost:8000/api/employees/ \
	-H "Content-Type: application/json" \
	-d '{"id":2,"name":"Yan","surname":"Fourth","phone":"44444444","company_id":1,"passport":{"type":"RF","number":"44444"},"department":{"name":"Develop2","phone":"22222"}}';
		curl -i -X POST localhost:8000/api/employees/ \
	-H "Content-Type: application/json" \
	-d '{"id":2,"name":"Masha","surname":"Fifth","phone":"555555","company_id":3,"passport":{"type":"RF","number":"55555"},"department":{"name":"Market","phone":"33333"}}';

# Запрос всех cотрудников 
get-AllEmployees:
	curl -i -X GET localhost:8000/api/employees/ \
	-H "Content-Type: application/json";

# Запрос всех cотрудников компании 1
get-AllEmployees1:
	curl -i -X GET localhost:8000/api/employees/?company_id=1 \
	-H "Content-Type: application/json";

# Запрос всех cотрудников компании 1 отдела Develop1
get-AllEmployees1Develop1:
	curl -i -X GET "localhost:8000/api/employees/?company_id=1&department_name=Develop1" \
	-H "Content-Type: application/json";

# Удалить сотрудников 1 и 2
delete-someEmployees:
	curl -i -X DELETE localhost:8000/api/employees/3 \
	-H "Content-Type: application/json";
	curl -i -X DELETE localhost:8000/api/employees/1 \
	-H "Content-Type: application/json";

# Поменять у 3 сотрудника телефон у 4 департамент
put-updateEmployee:
	curl -i -X PUT localhost:8000/api/employees/3 \
	-H "Content-Type: application/json" \
	-d '{"phone": "88005553535"}';
	curl -i -X PUT localhost:8000/api/employees/4 \
	-H "Content-Type: application/json" \
	-d '{"department":{"name":"NewDepart","phone":"7777777"}}';