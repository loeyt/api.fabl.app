FROM golang:1.16beta1 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./cmd/fabl

FROM gcr.io/distroless/base-debian10
COPY --from=build /app/fabl /
ENTRYPOINT [ "/fabl" ]
CMD [ "server" ]
