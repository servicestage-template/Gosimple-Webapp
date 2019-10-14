FROM golang:1.11

RUN mkdir -p /home/webapp/sample
ADD Gosimple-Webapp /home/webapp/sample/
ADD static /home/webapp/sample/static
ADD view /home/webapp/sample/view

WORKDIR /home/webapp/sample

RUN chmod +x /home/webapp/sample/Gosimple-Webapp
CMD ["/home/webapp/sample/Gosimple-Webapp"]

EXPOSE 8080
