
PASSWORD=qwerty
HOST_PORT=5432

CONTAINER_DB_NAME=employeesrestapi-db-1
CONTAINER_APP_NAME=employeesrestapi-employees-srv-1

all:

build: 
	docker-compose build employees-srv

run:
	docker-compose up employees-srv

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

test1: post-createLists post-createItems get-AllItems1
test2: put-doneItems get-AllItems1 get-AllItems1Done get-AllItems1DonePag

# Создать списки
post-createLists:
	curl -i -X POST localhost:8000/api/lists/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Покупка продуктов", "description": "На праздник"}';
	curl -i -X POST localhost:8000/api/lists/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Список домашнего задания", "description": "По курсам"}';
	curl -i -X POST localhost:8000/api/lists/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Встречи"}'


# Создать задания
post-createItems:
	curl -i -X POST localhost:8000/api/lists/1/items/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Хлеб", "description": "свежий"}';
	curl -i -X POST localhost:8000/api/lists/1/items/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Яблоки", "description": "2кг"}';
	curl -i -X POST localhost:8000/api/lists/1/items/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Молоко", "description": "3,2%"}';
	curl -i -X POST localhost:8000/api/lists/2/items/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Математика", "description": "2 задачи"}';
	curl -i -X POST localhost:8000/api/lists/2/items/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Информатика"}';
	curl -i -X POST localhost:8000/api/lists/3/items/ \
	-H "Content-Type: application/json" \
	-d '{"title": "Обсудить задачи", "description": "C отделом планирования"}';
	curl -i -X POST localhost:8000/api/lists/3/items/ \
	-H "Content-Type: application/json" \
	-d '{"title": "С другом"}'

# Отметить выполнено
put-doneItems:
	curl -i -X PUT localhost:8000/api/items/1 \
	-H "Content-Type: application/json" \
	-d '{"done": true, "description": "вчерашний"}';
	curl -i -X PUT localhost:8000/api/items/4 \
	-H "Content-Type: application/json" \
	-d '{"done": true}';
	curl -i -X PUT localhost:8000/api/items/7 \
	-H "Content-Type: application/json" \
	-d '{"done": true}';


# Удалить несколько задач
delete-Items:
	curl -i -X DELETE localhost:8000/api/items/2 \
	-H "Content-Type: application/json";
	curl -i -X DELETE localhost:8000/api/items/4 \
	-H "Content-Type: application/json";

# Запрос всех листа 1
get-AllItems1:
	curl -i -X GET localhost:8000/api/lists/1/items/ \
	-H "Content-Type: application/json";

# Запрос всех листа 1 с фильтрацией выполнено
get-AllItems1Done:
	curl -i -X GET localhost:8000/api/lists/1/items/?done=true \
	-H "Content-Type: application/json";

# Запрос всех листа 1 с фильтрацией не выполнено и пагинацией
get-AllItems1DonePag:
	curl -i -X GET "localhost:8000/api/lists/1/items/?done=false&limit=1&offset=0" \
	-H "Content-Type: application/json";
