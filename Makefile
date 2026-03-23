tree:
	tree . > .tree

run:
	clear
	go run .

start:
	clear
	docker compose down --remove-orphans
	CACHEBUST=$$(date +%s) docker compose build app
	docker compose up -d --force-recreate

stop:
	docker compose down --remove-orphans

logs:
	docker compose logs -f

.PHONY: tree run start stop logs