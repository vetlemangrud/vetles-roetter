FROM golang:1.27.1 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /out/data
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/vetles-roetter .

FROM gcr.io/distroless/static-debian12
WORKDIR /app
USER nonroot
COPY --from=build --chown=nonroot:nonroot /out /app
EXPOSE 8080
CMD ["/app/vetles-roetter"]
