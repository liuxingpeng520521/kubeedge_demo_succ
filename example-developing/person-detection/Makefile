.PHONY: all docker clean create delete

all:
	go build -o person-detection-app main.go
	docker build . -t mixiaochao/kubeedge-person-detection-app

docker:
	docker build . -t mixiaochao/kubeedge-person-detection-app

clean:
	rm -rf ./person-detection-app

delete:
	kubectl delete -f deployment.yaml

create:
	kubectl create -f deployment.yaml

push:
	docker push mixiaochao/kubeedge-person-detection-app