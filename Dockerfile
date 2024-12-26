FROM golang:latest as Builder

WORKDIR /fake-ai-detective
COPY . .
ENV CGO_ENABLED=0
RUN go build && chmod +x fake-ai-detective

FROM alpine:latest

COPY config/config.yaml /config/
COPY --from=Builder /fake-ai-detective/fake-ai-detective /

CMD ["/fake-ai-detective"]
