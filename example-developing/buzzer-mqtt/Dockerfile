FROM ubuntu:latest

CMD mkdir -p buzzer_mapper

COPY . buzzer_mapper/

WORKDIR buzzer_mapper

ENTRYPOINT ["/buzzer_mapper/main","-logtostderr=true"]
