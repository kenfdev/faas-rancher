TAG?=latest

build:
	docker build --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy -t kenfdev/faas-rancher:$(TAG) .

push:
	docker push kenfdev/faas-rancher:$(TAG)
