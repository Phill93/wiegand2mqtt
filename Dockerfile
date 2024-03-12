ARG BUILD_FROM

FROM golang:1.21
WORKDIR /src
COPY ./wiegand2mqtt /src
RUN CGO_ENABLED=0 go build -o wiegand2mqtt -tags rpi ./src/main.go 

FROM $BUILD_FROM
COPY --from=0 /src/wiegand2mqtt /bin/wiegand2mqtt
COPY ./run.sh /run.sh
CMD ["/run.sh"]
