FROM golang:1.19.2 AS builder
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /hello-app .

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=builder /hello-app /hello-app
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
CMD ["/hello-app"]
