FROM golang:1.23rc2

RUN apt-get update && \
  apt-get install -y --no-install-recommends mariadb-client-core qrencode && \
  apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /webapp/go
COPY ../webapp/go/go.* /webapp/go/
RUN go mod download

COPY ../webapp /webapp/
ENTRYPOINT [ "go", "run", "main.go" ]
