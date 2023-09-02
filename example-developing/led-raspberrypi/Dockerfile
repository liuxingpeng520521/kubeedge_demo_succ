FROM ubuntu:latest

CMD mkdir -p light_mapper

COPY . light_mapper/

WORKDIR light_mapper

ENTRYPOINT ["/light_mapper/main","-logtostderr=true"]
