FROM alpine:3.20

ARG TARGETARCH
COPY .build/linux-$TARGETARCH/minecraft_exporter /bin/minecraft_exporter

EXPOSE      9594
USER        nobody
ENTRYPOINT  [ "/bin/minecraft_exporter" ]
