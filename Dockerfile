FROM golang:1.24-bullseye AS build-go
ENV CGO_ENABLED=0

ARG BUILD_VERSION

WORKDIR /app
RUN go env -w GOMODCACHE=/root/.cache/go-build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . /app
RUN --mount=type=cache,target=/root/.cache/go-build go build -o overseer -ldflags "-X github.com/secmon-lab/overseer/pkg/domain/types.AppVersion=${BUILD_VERSION}" .

FROM gcr.io/distroless/base:nonroot
USER nonroot
COPY --from=build-go /app/overseer /overseer

ENTRYPOINT ["/overseer"]
