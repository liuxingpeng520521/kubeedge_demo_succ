.PHONY: all pi-execute-app docker clean

all: pi-execute-app

pi-execute-app:
	env GOOS=linux GOARCH=arm64 go build -o pi-counter-app main.go

docker:
	docker build . -t mixiaochao/kubeedge-pi-counter

clean:
	rm -f pi-counter-app

push:
	docker push mixiaochao/kubeedge-pi-counter

remove:
	rsync star pi@10.177.21.124:
