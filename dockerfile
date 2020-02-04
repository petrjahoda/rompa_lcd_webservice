FROM alpine:latest
RUN apk update && apk upgrade && apk add bash && apk add procps && apk add nano
RUN apk add samba-client
RUN apk add samba-common
RUN apk add cifs-utils
RUN apk add tzdata
RUN rm -rf /var/cache/apk/*
RUN cp /usr/share/zoneinfo/Europe/Prague /etc/localtime
WORKDIR /bin
COPY /css /bin/css
COPY /html /bin/html
COPY /js /bin/js
COPY /linux /bin
ENTRYPOINT rompa_lcd_webservice_linux