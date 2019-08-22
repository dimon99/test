FROM alpine:latest

ADD my_test_prom /app/custom_metric
WORKDIR /app
EXPOSE 8181
ENTRYPOINT ["/app/custom_metric"]
