FROM scratch
ADD bin/confmgr-static /confmgr-static
ADD confmgr.toml /confmgr.toml
EXPOSE 8080
ENTRYPOINT [ "/confmgr-static" ]
