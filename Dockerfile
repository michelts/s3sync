FROM golang:alpine
ADD src/ /opt/src/
WORKDIR /opt/src/
RUN go install
