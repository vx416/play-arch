
server:
	go run $(CURDIR)/main.go


build.image:
	docker build -f $(CURDIR)/dockerfile -t vicxu/server .
	docker push vicxu/server