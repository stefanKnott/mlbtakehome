FROM golang:1.20

WORKDIR /src/mlbtakehome
COPY . .
RUN go mod download
RUN GOOS=linux go build -o /mlbtakehome
EXPOSE 8080

CMD ["/mlbtakehome"]
