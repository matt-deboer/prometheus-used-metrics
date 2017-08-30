FROM alpine:3.5
COPY bin/prometheus-used-metrics /prometheus-used-metrics
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/prometheus-used-metrics"]
