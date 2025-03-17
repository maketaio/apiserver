FROM golang:1.24.0 AS dist
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o /apiserver ./cmd/apiserver

FROM gcr.io/distroless/base-debian11
COPY --from=dist /apiserver /apiserver
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/apiserver"]