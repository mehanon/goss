FROM golang:1.22-alpine3.19 AS env

WORKDIR /app
COPY shrimp/ /app/shrimp/
COPY cmd/ /app/cmd/
COPY go.mod go.sum /app/
RUN go build -o simp-shrimp ./cmd


FROM alpine:3.19

WORKDIR /app
RUN apk add --update --no-cache python3-dev py3-pip build-base ffmpeg &&\
    python3 -m pip install --break-system-packages cyberdrop-dl

COPY --from=env /app/simp-shrimp .

ENTRYPOINT ["./simp-shrimp"]
