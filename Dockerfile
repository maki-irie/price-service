FROM redhat/ubi9:latest AS BLD

RUN yum -y install go
# assume sources present in price-service.tar.gz in current directory
ADD price-service.tar.gz /home/eei-go-prj

ENV GOPATH=/home/eei-go-prj/price-service/cmd:/home/eei-go-prj/price-service/pkg/jwt:/home/eei-go-prj/price-service/pkg/postgres
RUN cd /home/eei-go-prj/price-service
RUN go build -ldflags="-s -w" -o ./cmd/service_run ./cmd/main.go

FROM redhat/ubi9-micro:latest
COPY --from=BLD /home/eei-go-prj/price-service/cmd/service_run /opt/app/

EXPOSE 8080
ENV REMOTE_SERVER=http://127.0.0.1:7070/
ENV DB_IP=http://127.0.0.1:5432

CMD /opt/app/service_run
