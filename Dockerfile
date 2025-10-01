FROM golang:1.24-alpine AS builder

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

RUN addgroup -g 65532 -S nonroot && adduser -u 65532 -S nonroot -G nonroot

COPY --from=builder /app/main /main

COPY --from=builder /app/env.yaml /home/nonroot/env.yaml

RUN chown -R nonroot:nonroot /home/nonroot

USER nonroot

WORKDIR /home/nonroot

EXPOSE 8080

ENTRYPOINT ["/main"]
