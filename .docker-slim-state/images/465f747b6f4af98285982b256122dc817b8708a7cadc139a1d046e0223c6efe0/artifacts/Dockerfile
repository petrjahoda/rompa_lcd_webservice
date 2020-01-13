FROM scratch
LABEL docker-slim.version="darwin|Transformer|1.26.1|2ec04e169b12a87c5286aa09ef44eac1cea2c7a1|2019-11-28_04:37:59PM"
ADD files.tar /
WORKDIR /bin
ENV PATH "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
ENTRYPOINT ["/bin/sh","-c","rompa_lcd_webservice_linux"]
