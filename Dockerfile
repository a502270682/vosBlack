FROM  golang:latest

#ENV GO111MODULE on
#WORKDIR /test/cache
#ADD go.mod .
#ADD go.sum .
#RUN go mod download
#WORKDIR /test/release
#ADD build .
#RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o vosBlack main.go

WORKDIR /root
COPY /vosBlack /
COPY /config/vosBlack.toml /etc/config/vosBlack.toml

EXPOSE 8080
CMD ["/vosBlack", "-c", "/etc/config/vosBlack.toml"]
