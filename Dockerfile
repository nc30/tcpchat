FROM golang:alpine
ENV APP_ROOT=/opt
RUN mkdir -p $APP_ROOT
WORKDIR $APP_ROOT
EXPOSE 23

ADD . $APP_ROOT
CMD ["go", "run", "main.go"]
