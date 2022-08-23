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
RUN CGO_ENABLED=0 go build -o signer http/main.go

#
# Final Stage
#

FROM alpine:3.16

WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /src/build/signer .
RUN chmod +x signer
EXPOSE 8080
CMD [ "./signer" ]