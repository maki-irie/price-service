FROM redhat/ubi9-micro:latest
ADD ubi9-bin/cmd /opt/app/
ENV GOPATH=/opt/app/

EXPOSE 8080
ENV REMOTE_SERVER=http://127.0.0.1:7070/
ENV DB_IP=http://127.0.0.1:5432

CMD /opt/app/cmd
