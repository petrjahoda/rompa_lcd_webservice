#!/usr/bin/env bash
cd linux
upx rompa_lcd_webservice_linux
cd ..
docker rmi -f petrjahoda/rompa_lcd_webservice:latest
docker build -t petrjahoda/rompa_lcd_webservice:latest .
docker push petrjahoda/rompa_lcd_webservice:latest