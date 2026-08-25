FROM golang:1.27-bookworm
ENV GOPROXY=off GOSUMDB=off GOFLAGS=-mod=vendor
WORKDIR /src
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /out/metricsd ./cmd/metricsd
CMD ["/out/metricsd"]
