FROM sammobach/go:1.16 as build
LABEL maintainer="Sam Mobach <sam@hoi.studio>"
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -ldflags "-s -w" -o dist/secretary ./cmd/secretary
CMD ["/app/dist/secretary"]

FROM alpine:3.14.0
LABEL maintainer="Sam Mobach <sam@hoi.studio>"
RUN mkdir /app
COPY --from=build /app /app
WORKDIR /app
CMD ["/app/dist/secretary"]
