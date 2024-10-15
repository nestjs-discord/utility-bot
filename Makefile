.PHONY: redeploy
redeploy:
	docker compose down --remove-orphans
	docker compose --env-file ./.env up -d --build
	sleep 5
	docker compose ps
	docker stats --no-stream
	docker compose logs
