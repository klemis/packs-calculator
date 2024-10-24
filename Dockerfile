FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ../.. .

# Build the packs-calculator-service API binary.
RUN go build -o api ./cmd

# Set the default command to run the API binary.
CMD ["./api"]