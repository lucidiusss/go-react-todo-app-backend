# build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy files
COPY . .

# compile app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# final
FROM alpine:latest

# sertificates
RUN apk --no-cache add ca-certificates

# create user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# switch user
USER appuser

WORKDIR /home/appuser/

# copy app file
COPY --from=builder /app/main .

# listen port
EXPOSE 3000

# start the app
CMD ["./main"]
