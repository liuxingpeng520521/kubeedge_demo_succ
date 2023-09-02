
# make led_light_mapper
.PHONY: default led_light_mapper push delete create
led_light_mapper:
	env GOOS=linux GOARCH=arm64 go build main.go
	docker build -t mixiaochao/led-light-mapper .

push:
	docker push mixiaochao/led-light-mapper

move:
	rsync  main pi@10.177.21.124:/home/pi

delete:
	kubectl delete -f deployment.yaml

create:
	kubectl create -f deployment.yaml
