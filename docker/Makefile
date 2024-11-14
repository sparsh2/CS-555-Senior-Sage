# docker user name
NAMESPACE := sarodesparsh

all: authz-docker storage-service-docker

authz-docker:
	docker build -t $(NAMESPACE)/authz -f ./authz/Dockerfile ../backend/authz/
	docker push $(NAMESPACE)/authz

storage-service-docker:
	docker build -t $(NAMESPACE)/storage-service -f ./storage-service/Dockerfile ../backend/storage-service/
	docker push $(NAMESPACE)/storage-service

authz-test-docker:
	docker build -t $(NAMESPACE)/test-authz -f ./test-authz/Dockerfile ../backend/authz/
	docker run --rm $(NAMESPACE)/test-authz
	docker rmi $(NAMESPACE)/test-authz

storage-service-test-docker:
	docker build -t $(NAMESPACE)/storage-service-test -f ./storage-service-test/Dockerfile ../backend/storage-service/
	docker run --rm $(NAMESPACE)/storage-service-test

test: authz-test-docker
	

.PHONY: all authz-docker storage-service-docker