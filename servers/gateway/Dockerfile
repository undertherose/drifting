FROM alpine
RUN apk add --no-cache ca-certificates
COPY finalgateway /finalgateway
EXPOSE 443
ENTRYPOINT ["/finalgateway"]