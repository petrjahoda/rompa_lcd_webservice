#!/usr/bin/env bash
#cd linux
#upx rompa_lcd_webservice_linux
#cd ..
#docker rmi -f petrjahoda/rompa_lcd_webservice:latest
#docker build -t petrjahoda/rompa_lcd_webservice:latest .
#docker push petrjahoda/rompa_lcd_webservice:latest


./update
name=${PWD##*/}
go get -u all
GOOS=linux go build -ldflags="-s -w" -o linux/"$name"
cd linux
upx "$name"
cd ..

docker rmi -f petrjahoda/"$name":latest
docker  build -t petrjahoda/"$name":latest .
docker push petrjahoda/"$name":latest

docker rmi -f petrjahoda/"$name":2020.4.2
docker build -t petrjahoda/"$name":2020.4.2 .
docker push petrjahoda/"$name":2020.4.2