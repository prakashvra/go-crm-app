FROM golang:1.17-alpine AS build
WORKDIR /src
COPY . .
ENV CGO_ENABLED=0
RUN go build -o ../bin/crm-app src/server/main.go
FROM scratch AS bin
COPY --from=build ../bin/crm-app /
ENTRYPOINT ["/crm-app"]
