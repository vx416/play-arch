
consul.devagent:
	consul agent -dev -enable-script-checks -config-dir=$(CURDIR)/consul.d

consul.reload:
	consul reload

server_a:
	go run $(CURDIR)/services/server-a/main.go -port=$(p) -app_name=$(n)

server_b:
	go run $(CURDIR)/services/server-b/main.go -port=$(p) -app_name=$(n)

consul.sidecar:
	consul connect proxy -sidecar-for $(svc)

consul.sidecar.envoy:
	consul connect envoy -sidecar-for $(svc) -admin-bind localhost:$(admin_port)

consul.health:
	curl 'http://localhost:8500/v1/health/service/$(svc)?passing=1&connect=true'