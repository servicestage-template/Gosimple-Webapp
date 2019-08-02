FROM frolvlad/alpine-glibc:latest

RUN mkdir -p /home/webapp/sample
ADD main /home/webapp/sample/
ADD static /home/webapp/sample/static
ADD view /home/webapp/sample/view

WORKDIR /home/webapp/sample

RUN chmod +x /home/webapp/sample/main
CMD ["/home/webapp/sample/main"]

EXPOSE 8080
