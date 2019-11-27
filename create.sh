#!/usr/bin/env bash
docker rmi -f petrjahoda/rompa_lcd_webservice:"$1"
docker build -t petrjahoda/rompa_lcd_webservice:"$1" .
docker push petrjahoda/rompa_lcd_webservice:"$1"