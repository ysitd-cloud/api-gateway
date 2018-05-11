FROM ysitd/dep AS builder

WORKDIR /go/src/app.ysitd/gateway

COPY . .

RUN dep ensure -v -vendor-only && \
    go install -v

FROM ysitd/binary

COPY --from=builder /go/bin/gateway /gateway

CMD ["/gateway"]
