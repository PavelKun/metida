APP = metida
VER = 0.0.1

#not working?
.PHONY: build
build:
	go build -o $(APP)

#not working?
.PHONY: build-image
build-image:
	docker build -t $(APP):$(VER) .

#not working?
.PHONY: run-app
run-app:
	docker run -d --name=$(APP)-$(VER) -p 8080:8080 -v /home/.../metida/db:/go/db $(APP):$(VER)

.PHONY: del-app
del-app:
	docker rm $(APP)-$(VER)


.PHONY: swag-init
swag-init:
	swag init

# пока служит для запуска только бд
.PHONY: db-start
db-start:
	docker-compose up -d

# кодогенерация из sql скриптов
# https://docs.sqlc.dev/en/latest/index.html
.PHONY: sql-generate
sql-generate:
	sqlc generate


