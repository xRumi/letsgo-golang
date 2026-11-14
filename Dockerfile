FROM golang:1.27 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN mkdir build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/ ./cmd/web/

# scratch build copying just the build files
FROM scratch

WORKDIR /app
COPY --from=build /app/build/ ./build/
COPY --from=build /app/ui/ ./ui/

EXPOSE 4000
CMD ["./build/web"]