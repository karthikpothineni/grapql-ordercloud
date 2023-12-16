.PHONY:docker-build docker-run docker-start docker-stop

GONAME=product-service

default: docker-build

docker-build:
	@docker build --tag $(GONAME) .

docker-run:
	@docker run --env-file ./.env -p 8080:8080 --name $(GONAME) $(GONAME)

docker-stop:
	@docker stop $(GONAME)

docker-start:
	@docker start $(GONAME)