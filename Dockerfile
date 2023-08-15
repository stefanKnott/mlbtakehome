FROM golang:1.20

WORKDIR /src/mlbtakehome
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux go build -o /mlbtakehome
EXPOSE 8080

CMD ["/mlbtakehome"]
