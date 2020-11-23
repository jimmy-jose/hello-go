FROM golang:1.15.3
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go get -d github.com/gorilla/mux
RUN go build -o main .
CMD ["/app/main"]