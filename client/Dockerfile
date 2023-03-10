FROM golang:1-alpine AS build_client
WORKDIR /app
COPY . .
RUN go build -o client .

FROM alpine:3
COPY --from=build_client /app/client /app/client
CMD ["/app/client"]