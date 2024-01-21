FROM golang:1.20.1-alpine3.17 AS builder
COPY . /build
WORKDIR /build
RUN apk add make git 
RUN apk add --no-cache ca-certificates
RUN make dependencies
RUN make
WORKDIR /permissions
RUN echo "carbide-images-api:x:1001:1001::/:" > passwd && echo "carbide-images-api:x:2000:carbide-images-api" > group

FROM scratch
COPY --from=builder /permissions/passwd /etc/passwd
COPY --from=builder /permissions/group /etc/group
USER carbide-images-api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/carbide-images-api /
ENTRYPOINT ["/carbide-images-api"]
