FROM golang:1.22-bullseye AS build-go
ENV CGO_ENABLED=0

ARG BUILD_VERSION
COPY . /app
WORKDIR /app
RUN go build -o overseer -ldflags "-X github.com/m-mizutani/overseer/pkg/domain/types.AppVersion=${BUILD_VERSION}" .

FROM gcr.io/distroless/base
COPY --from=build-go /app/overseer /overseer

ENTRYPOINT ["/overseer"]
