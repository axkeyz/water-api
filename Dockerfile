# syntax=docker/dockerfile:1

# BUILD
FROM golang:1.17.3-buster AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /water

# DEPLOY
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /water /water

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/water"]