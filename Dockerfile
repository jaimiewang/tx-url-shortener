FROM golang:1.14 as build
WORKDIR /build
COPY . /build
RUN cd src/ && go build -o ../bin/tx-url-shortener .

FROM debian:stable
RUN useradd tx-url-shortener
WORKDIR /opt/tx-url-shortener/bin
COPY --from=build /build/bin/tx-url-shortener .
WORKDIR ..
RUN chown -hR tx-url-shortener:tx-url-shortener /opt/tx-url-shortener
USER tx-url-shortener
EXPOSE 8080
ENTRYPOINT ["./bin/tx-url-shortener"]