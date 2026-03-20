FROM golang:1.26.1

RUN mkdir /app
WORKDIR /app

RUN apt-get update && apt-get upgrade -y && apt-get install -y make git

RUN go install -v github.com/cespare/reflex@latest
ENTRYPOINT ["reflex", "-c", "./reflex.conf"]
