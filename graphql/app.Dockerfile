FROM ubuntu:latest
LABEL authors="mstare"

ENTRYPOINT ["top", "-b"]