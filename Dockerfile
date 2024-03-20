FROM golang:1.22.1-alpine AS modules
WORKDIR /app
COPY go.mod go.sum ./
RUN  go mod download
COPY . ./

FROM modules AS server-build
RUN go build -o server ./cmd/server/main.go

FROM modules AS client-build
RUN go build -o client ./cmd/client/main.go

FROM scratch AS server
COPY --from=server-build /app/server /app/.env ./
CMD ["./server"]

FROM scratch AS client
COPY --from=client-build /app/client /app/.env ./
CMD ["./client"]
