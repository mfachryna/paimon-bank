# run app
.PHONY: run
run:
	@go run cmd/main.go

# build app
.PHONY: build
build:
	@go build -o ./build/app ./cmd/main.go

# build app alpine
.PHONY: build-alpine
build-alpine:
	@go mod tidy && \
	GOOS=linux GOARCH=amd64 go build -o ./build/app ./cmd/main.go

# migrate up
.PHONY: migrate-up
migrate-up:
	migrate -path db/migrations -verbose \
	-database "postgres://postgres:root@localhost:5432/paimon_bank" up

# migrate down
.PHONY: migrate-down
migrate-down:
	migrate -path db/migrations -verbose \
	-database "postgres://postgres:root@localhost:5432/paimon_bank" down

# make startProm
.PHONY: start-prom
start-prom:
	docker run \
	--rm \
	--network="host" \
	-p 9090:9090 \
	--name=prometheus \
	-v $(shell pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
	prom/prometheus

# make startGrafana
# for first timers, the username & password is both `admin`
.PHONY: start-grafana
start-grafana:
	docker volume create grafana-storage
	docker volume inspect grafana-storage
	docker run -p 3000:3000 --name=grafana grafana/grafana-oss || docker start grafana
	