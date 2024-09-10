FROM redhat/ubi9:latest AS BLD

RUN yum -y install go

RUN mkdir -p /src/price-service
WORKDIR /src/price-service
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -ldflags="-s -w" -o price-service cmd/main.go

FROM redhat/ubi9-micro:latest

COPY --from=BLD /src/price-service/price-service /appl/

EXPOSE 8080
ENV REMOTE_SERVER=http://127.0.0.1:7070/
ENV DB_IP=http://127.0.0.1:5432

WORKDIR /appl/
CMD ["./price-service"]
