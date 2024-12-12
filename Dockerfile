FROM registry.gitlab.com/invx/devops/templates/images/golang:1.23-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache tzdata

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build main.go

FROM registry.gitlab.com/invx/devops/templates/images/alpine:3.20
WORKDIR /app

COPY --from=builder /build/main .
COPY .env .
ENV TZ=Asia/Bangkok
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
EXPOSE 8910

CMD ["/app/main"]