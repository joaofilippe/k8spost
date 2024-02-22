FROM golang as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download



COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./server/server ./server/main.go

FROM scratch
COPY --from=builder /app/server/server /app/server/server

ENV ENVI=docker

WORKDIR /app/server
ENTRYPOINT [ "./server" ]