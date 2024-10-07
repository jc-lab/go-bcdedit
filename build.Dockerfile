FROM golang:1.22-alpine3.20 as builder

RUN apk add \
    ca-certificates git ocaml gcc make automake autoconf pkgconfig m4 libtool \
    tar xz gettext-dev \
    musl-dev

RUN mkdir -p /build && \
    cd /build && \
    git clone https://github.com/libguestfs/hivex.git && \
    cd hivex && \
    git checkout -f 7996d79cf910758fcb29554e83f44c94560f170d && \
    autoreconf -i && \
    ./generator/generator.ml && \
    ./configure && \
    make && \
    make check || (cat config.log && false) && \
    make install

RUN mkdir -p /build/src /build/dist
COPY . /build/src

WORKDIR /build/src
RUN go build -o /build/dist/bcdedit-linux --ldflags '-linkmode external -extldflags "-static"' ./cmd/bcdedit/main.go

FROM scratch
COPY --from=builder /build/dist/ /