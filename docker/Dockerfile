FROM golang:1.26.4-trixie AS constructor

WORKDIR /app
COPY ./src ./src
COPY go.mod go.sum ./
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/gt ./src

FROM golang:1.26.4-trixie AS runtime-stage

WORKDIR /
COPY --from=constructor /app/gt /gt
COPY .env /.env
COPY item-constants.json /item-constants.json
COPY unusuals.json /unusuals.json
ENTRYPOINT ["/gt"]
