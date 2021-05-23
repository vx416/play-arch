

.PHONY: gen.proto

gen.proto:
	protoc -I=. --go_out=plugins=grpc,paths=source_relative:. ./proto/*.proto

run.jaeger:
	docker run -d --name jaeger \
	-e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
	-p 5775:5775/udp \
	-p 6831:6831/udp \
	-p 6832:6832/udp \
	-p 5778:5778 \
	-p 16686:16686 \
	-p 14268:14268 \
	-p 14250:14250 \
	-p 9411:9411 \
	jaegertracing/all-in-one:1.22

run.kafka:
	docker-compose -f $(PWD)/kafka/docker-compose.yaml -p kafka up

machine.up:
	vagrant up

machine.down:
	vagrant down