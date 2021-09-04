# Build Stage
FROM golang:1.17.0-bullseye AS build-stage

WORKDIR /app

ADD . .

RUN make lint test build

# Final Stage
FROM debian:buster-slim

WORKDIR /app

COPY --from=build-stage /app/bin/proxy bin/proxy/proxy
RUN chmod +x bin/proxy/proxy

COPY --from=build-stage /app/proxy-configs/ proxy-configs/
RUN chmod +r proxy-configs/

CMD ["bin/proxy/proxy"]
