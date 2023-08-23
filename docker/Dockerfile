## Builder Layer
FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /k8s-from-secrets-vault

## Deployed Layer
FROM alpine AS run

WORKDIR /

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /k8s-from-secrets-vault /k8s-from-secrets-vault

## Golang runtime dependency
RUN apk add --no-cache gcompat

## Add executable permissions to app file owners group
RUN chmod -R 777 /k8s-from-secrets-vault
## Make application-runners group owner of app binaries file
RUN chown -R root:root /k8s-from-secrets-vault

CMD ["/k8s-from-secrets-vault"]