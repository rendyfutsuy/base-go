FROM golang:1.22.5-bookworm as builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
COPY . ./
RUN go build main.go

FROM debian:bookworm as production
RUN apt update
RUN apt install -y curl ca-certificates
COPY --from=builder app .
EXPOSE 9000
CMD ./main