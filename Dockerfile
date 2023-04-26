FROM golang:alpine AS builder
LABEL authors="rysiw"
WORKDIR /app
COPY . /app
ENV GOPROXY https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch AS runner
WORKDIR /app
COPY  --from=builder /app/ .
EXPOSE 8000
ENTRYPOINT ["./main"]