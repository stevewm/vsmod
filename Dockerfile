FROM debian:bookworm-slim

COPY vsmod /

ENTRYPOINT ["/vsmod"]
CMD ["download", "--file", "/config/mods.yaml"]
