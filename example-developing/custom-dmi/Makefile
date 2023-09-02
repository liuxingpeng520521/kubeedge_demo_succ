# make modbus-dmi
.PHONY: default build push delete create deviceapply devicedelete
build:
	env GOOS=linux GOARCH=arm64 go build -o modbus-dmi main.go

push:
	docker push mixiaochao/led-light-mapper

move:
	rsync  main config.yaml pi@10.177.21.124:/home/pi/hd/kubeedge-example/

delete:
	kubectl delete -f deployment.yaml

create:
	kubectl create -f deployment.yaml
deviceapply:
	kubectl apply -f resource/random-device-model.yaml
	kubectl apply -f resource/random-device-instance.yaml
devicedelete:
	kubectl delete -f resource/random-device-instance.yaml
	kubectl delete -f resource/random-device-model.yaml

run:
	./main --v 4 --config-file=../config.yaml