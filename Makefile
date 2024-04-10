.PHONY: run

all: pub

pub:
	cd stan-pub && go run ./cmd/publisher.go
sub:
	cd stan-sub && go run ./cmd/subscriber.go
nats:
	sudo docker run -d --name nats-streaming -p 4222:4222 nats-streaming:latest