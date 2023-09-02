# make modbus-dmi
.PHONY: default build push delete create deviceapply devicedelete
build:
	env GOOS=linux GOARCH=arm64 go build -o modbus-dmi main.go

push:
	docker push mixiaochao/led-light-mapper

move:
	rsync  modbus-dmi pi@10.177.21.124:/home/pi/hd/kubeedge-example/modbus-dmi

delete:
	kubectl delete -f deployment.yaml

create:
	kubectl create -f deployment.yaml
deviceapply:
	kubectl apply -f model-modbus.yaml
	kubectl apply -f modbus.yaml
devicedelete:
	kubectl delete -f modbus.yaml
	kubectl delete -f model-modbus.yaml