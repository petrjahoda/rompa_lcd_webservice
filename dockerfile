FROM alpine:latest
RUN apk update && apk upgrade && apk add bash && apk add procps && apk add nano
RUN apk add samba-client
RUN apk add samba-common
RUN apk add cifs-utils
RUN rm -rf /var/cache/apk/*
WORKDIR /bin
COPY /linux /bin
COPY /css /bin/css
COPY /html /bin/html
COPY /js /bin/js
ENTRYPOINT rompa_lcd_webservice_linux