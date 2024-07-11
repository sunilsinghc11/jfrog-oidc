JOIN_KEY=e405527b8f63f94a85d363211efbf4b9

build:
	go build -o oidc-poc github.com/mosheya/access-oidc-poc/oidc-service

run:
	go run ./oidc-service/server.go

start-artifactory:
	docker run --rm -d --name arti --rm -p 8082:8082 \
  -v $(CURDIR)/artifactory:/var/opt/jfrog/artifactory \
  -e JF_SHARED_SECURITY_JOINKEY=${JOIN_KEY} \
  releases-docker.jfrog.io/jfrog/artifactory-pro:7.59.11

stop-artifactory:
	docker stop arti

create-token:
	$(jfrog access st --url=http://localhost:8082/artifactory --join-key=e405527b8f63f94a85d363211efbf4b9)

ngrok:
	ngrok http 8080

ngrok-nginx: # ngrok support only one session
	ngrok http 8085

start-nginx:
	docker run -d -p 8085:8085 --name nginx -v $(CURDIR)/nginx-oidc.conf:/etc/nginx/conf.d/default.conf:ro nginx:1.25.1
