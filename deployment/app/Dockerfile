FROM golang:1.21.3 as build

WORKDIR /app

COPY . .
RUN go mod download

ARG BUILD_VERSION
ARG VARIABLES_PATH

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/app  \
    -ldflags "-X $VARIABLES_PATH.release=$BUILD_VERSION -X $VARIABLES_PATH.buildDate=$(date -u +%Y-%m-%dT%H:%M:%S) -X $VARIABLES_PATH.gitHash=$(git log --format="%h" -n 1)" \
    ./cmd/app

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/app ./bin/

COPY ./config/ ./config/

ENTRYPOINT ["./bin/app"]