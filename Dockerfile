FROM golang as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./server/server ./server/main.go


ENV ENVI=docker

ENTRYPOINT [ "./server/server" ]