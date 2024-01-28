FROM golang:1.21 as build

WORKDIR /usr/src/otus-socialmedia

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o otus-socialmedia-api

FROM alpine:3.11

COPY --from=build /usr/src/otus-socialmedia/otus-socialmedia-api otus-socialmedia-api

EXPOSE 4000

CMD [ "./otus-socialmedia-api" ]
