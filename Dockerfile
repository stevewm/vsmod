FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
    
COPY vsmod /

ENTRYPOINT ["/vsmod"]
CMD ["download", "--file", "/config/mods.yaml"]
