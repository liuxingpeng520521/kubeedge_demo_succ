
# make led_light_mapper
.PHONY: default buzzer_mqtt_mapper push delete create
buzzer_mqtt_mapper:
	env GOOS=linux GOARCH=arm64 go build main.go
	docker build -t mixiaochao/buzzer-mqtt-mapper .

push:
	docker push mixiaochao/buzzer-mqtt-mapper

move:
	rsync  main mixiaochao@10.177.29.149:/home/mixiaochao

delete:
	kubectl delete -f deployment.yaml

create:
	kubectl create -f deployment.yaml
