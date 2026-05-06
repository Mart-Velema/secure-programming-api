FROM golang:1.26.2-trixie AS constructor

WORKDIR /app
COPY ./src ./
COPY go.mod go.sum ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /sp

FROM golang:1.26.2-trixie AS runtime-stage

WORKDIR /
COPY --from=constructor /sp /sp
COPY .env /.env
ENTRYPOINT ["/sp"]
