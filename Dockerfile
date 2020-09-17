FROM golang:1.14-alpine as builder

WORKDIR /go/src/cpid-solar-telemetry
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
 
RUN go build

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/cpid-solar-telemetry .

CMD ["./cpid-solar-telemetry"]
