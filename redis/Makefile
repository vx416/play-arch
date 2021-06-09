
PWD ?= $(pwd)

run.single:
	docker run --name redis -d -p 6379:6379 \
	-v $(PWD)/redis.conf:/etc/redis.conf \
	redis redis-server /etc/redis.conf

run.cluster:
	docker-compose -f $(PWD)/docker-compose.yaml -p redis-cluster up --build