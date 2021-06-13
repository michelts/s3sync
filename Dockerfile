FROM golang:alpine
WORKDIR /opt/src/
ADD src/ /opt/src/
#RUN go install
