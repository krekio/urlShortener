FROM ubuntu:latest
LABEL authors="Andrew"

ENTRYPOINT ["top", "-b"]