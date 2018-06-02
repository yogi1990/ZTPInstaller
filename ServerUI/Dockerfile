FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/test/PNP

ENV USER ykumar
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET 00MTi32Qr4BM1u5D

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://ykumar@localhost:5432/PNP?sslmode=disable

WORKDIR /go/src/github.com/test/PNP

RUN godep go build

EXPOSE 8888
CMD ./PNP