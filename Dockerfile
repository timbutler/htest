FROM golang
RUN  apt-get update \
  && apt-get install -y wget \
  && rm -rf /var/lib/apt/lists/*
RUN wget --no-check-certificate -O /htest https://github.com/timbutler/htest/releases/download/1.1/htest
RUN chmod +x /htest
CMD ["/htest"]

EXPOSE 8000
