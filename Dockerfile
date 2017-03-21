FROM golang
COPY htest /
CMD ["/htest"]

EXPOSE 8000
