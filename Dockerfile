FROM golang AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /takehome ./cmd

# final container
FROM alpine
COPY --from=builder /takehome /takehome
ENTRYPOINT ["/takehome"]
