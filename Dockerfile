FROM golang:alpine

#RUN apk add curl

COPY src /client

WORKDIR /client

# Run just the first time; then comment the line:
RUN go mod init myclient && go mod tidy

#RUN go mod tidy
#CMD go run /client/client.go

CMD sleep infinity