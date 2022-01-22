FROM moby/buildkit:v0.9.3
WORKDIR /render
COPY render README.md /render/
ENV PATH=/render:$PATH
ENTRYPOINT [ "/bhojpur/render" ]