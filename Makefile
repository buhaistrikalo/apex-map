build: 
	docker-compose build
up: 
	docker-compose up -d
down: 
	docker-compose down
restart:
	docker-compose down
	docker-compose build
	docker-compose up -d