FROM golang:1.18 as builder

LABEL maintainer="kiettiphong.m@avareum.finance"

#
# Building
#

WORKDIR /src/build
COPY go.mod /src/build/
COPY go.sum /src/build/
RUN go mod download

ADD . /src/build
RUN go build -o signer_app http/main.go

#
# Final Stage
#

FROM alpine:3.16

WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /src/build/signer_app .
EXPOSE 8080
CMD ["./signer_app"]