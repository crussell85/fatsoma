FROM golang:1.22 as build

WORKDIR /src

COPY src ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
COPY --from=build /src/ticket-api /ticket-api
EXPOSE 3000
ENTRYPOINT ["/ticket-api"]