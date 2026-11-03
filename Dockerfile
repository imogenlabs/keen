# Used by goreleaser (dockers_v2): it copies the binaries it already built.
# distroless/static has CA certificates and runs as a non-root user.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/keen $TARGETPLATFORM/keen-mcp /usr/bin/
# keen-mcp by default; `docker run <image> keen ...` runs the CLI instead.
CMD ["keen-mcp"]
