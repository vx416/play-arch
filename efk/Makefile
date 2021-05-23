

run:
	go run ./server/*.go

build.image:
	docker build -f ./server/Dockerfile	-t vicxu/pingpong .

run.es:
	docker run -d -p 9200:9200 -p 9300:9300 --name es \
  -e "http.cors.enabled=true" -e "http.cors.allow-origin=http://localhost:8080" -e "discovery.type=single-node" \
  docker.elastic.co/elasticsearch/elasticsearch:7.9.1

run.app:
	docker run -d -p 13333:13333 --name app \
	-e LOGFILE=/tmp/log -v $(CURDIR)/tmp/log:/tmp/log vicxu/pingpong:latest

run.fld:
	docker run -d -p 24224:24224 -p 24224:24224/udp \
	-v 	$(CURDIR)/collector/fluentd.conf:/fluentd/etc/fluent.conf -v /tmp/data:/fluentd/log \
	-e FLUENT_CONF=/fluentd/etc/fluent.conf --name fluentd fluent/fluentd:v1.3-debian-1

run.flb:
	docker run -p 127.0.0.1:24224:24224 -d  --name flb \
	-v $(CURDIR)/collector/fluentbit.conf:/etc/fluentbit.conf \
	-v $(CURDIR)/tmp/log:/tmp/log \
	fluent/fluent-bit:1.5 /fluent-bit/bin/fluent-bit -c /etc/fluentbit.conf


