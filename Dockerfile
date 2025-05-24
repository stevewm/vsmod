FROM gcr.io/distroless/static-debian12

COPY vsmod /

ENTRYPOINT ["/vsmod"]
CMD ["download", "--file", "/config/mods.yaml"]
