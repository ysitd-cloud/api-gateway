FROM ysitd/dep AS builder

WORKDIR /go/src/code.ysitd.cloud/gateway

COPY . .

RUN dep ensure -v -vendor-only && \
    go install -v

FROM ysitd/binary

COPY --from=builder /go/bin/gateway /gateway

CMD ["/gateway"]
