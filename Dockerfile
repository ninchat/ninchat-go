FROM ubuntu

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get -y --no-install-recommends install chromium-browser firefox && \
    apt-get clean && \
    adduser --disabled-password --disabled-login --gecos "" user

USER user
