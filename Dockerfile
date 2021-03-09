FROM scratch
EXPOSE 5473/tcp
COPY nettest /nettest
COPY pki/ca.crt /tmp/ca.crt
COPY pki/client.crt /tmp/client.crt
COPY pki/client.key /tmp/client.key
ENTRYPOINT ["/nettest"]
