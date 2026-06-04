FROM golang:1.22 as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /
COPY --from=builder /app/app /app

EXPOSE 8080
CMD ["/app"]
