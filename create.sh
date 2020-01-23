#!/usr/bin/env bash
docker rmi -f petrjahoda/rompa_lcd_webservice:latest
docker build -t petrjahoda/rompa_lcd_webservice:latest .
docker push petrjahoda/rompa_lcd_webservice:latest