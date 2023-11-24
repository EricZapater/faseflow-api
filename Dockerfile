FROM golang:alpine AS build

RUN apk add --update git
WORKDIR /go/src/github.com/EricZapater/faseflow-api
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/faseflow-api cmd/api/main.go

FROM scratch
COPY --from=build /go/bin/faseflow-api /go/bin/faseflow-api
ENTRYPOINT ["/go/bin/faseflow-api"]