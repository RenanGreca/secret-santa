FROM --platform=$BUILDPLATFORM golang:alpine AS build

ARG VERSION
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN apk add build-base

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o secret-santa ./

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/secret-santa .
COPY ./players/players.txt ./players/players.txt

RUN /app/secret-santa shuffle

EXPOSE 3000

ENTRYPOINT ["/app/secret-santa"]
