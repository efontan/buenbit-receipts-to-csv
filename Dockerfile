FROM golang:1.16-alpine AS builder

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM ubuntu:18.04
COPY --from=builder ./app .

# install sajari/docconv dependencies
RUN apt-get update && \
     apt-get install -y poppler-utils wv unrtf tidy

CMD ["./app"]