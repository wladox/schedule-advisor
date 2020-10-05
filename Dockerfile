FROM golang:1.14-stretch AS builder

ARG VERSION
ENV VERSION ${VERSION}

# Make app
ENV APP_DIR /Users/wladimir.gaus/workspace/schedule-advisor
RUN mkdir -p $APP_DIR
COPY . $APP_DIR
WORKDIR $APP_DIR

# install app dependencies
RUN go build

# Add Certs
FROM ubuntu:bionic
RUN apt-get update && TERM=linux DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
  && rm -rf /var/lib/apt/lists/*

# Copy complied binary
COPY --from=builder /Users/wladimir.gaus/workspace/schedule-advisor/out/schedule-advisor .

#Run app and expose port
ARG PORT=8000
ENV PORT=${PORT}

ENV FLINK_METRICS_API ""
ENV LOG_LEVEL "DEBUG"

CMD ["./schedule-advisor"]

EXPOSE ${PORT}
