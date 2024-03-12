FROM golang:1.21
WORKDIR /src
COPY ./wiegand2mqtt /src
RUN go build -o /bin/wiegand2mqtt -tags rpi ./src/main.go 

FROM alpine:latest
COPY --from=0 /bin/wiegand2mqtt /bin/wiegand2mqtt
COPY ./wiegand2mqtt/config.yaml.example .
CMD ["/bin/wiegand2mqtt"]
