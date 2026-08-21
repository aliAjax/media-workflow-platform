FROM golang:1.26 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/media-api ./cmd/media-api
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/media-api /media-api
EXPOSE 8084
ENTRYPOINT ["/media-api"]
