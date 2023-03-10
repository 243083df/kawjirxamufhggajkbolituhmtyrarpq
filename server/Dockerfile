FROM golang:1-alpine AS build_server
WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:3
COPY --from=build_server /app/server /app/server
CMD ["/app/server"]