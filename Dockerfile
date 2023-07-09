#FROM указываем какой образ докер использовать для сборки 
#AS BUILDER -присваеваем название файлу
FROM golang:1.19 AS BUILDER 
#показываем версию GO
RUN go version
#устанавливаем пакет git
#копируем в каталог /github
COPY ./ /github/
#используем каталог github как рабочий
WORKDIR /github/
RUN pwd
RUN ls -la
#загружаем необходимые для сборки зависимости
RUN go mod download && go get -u ./...
#запускаем GO сборку
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/viewvrn ./main.go

#используем более легковесный docker образ
FROM ubuntu:latest
#установка необходимых пакетов
RUN apt-get update && apt-get install git cron -y
RUN mkdir /data
#используем директорию root как рабочую
WORKDIR /github/
RUN mkdir log
RUN mkdir templates
RUN pwd
RUN ls -la
#копируем в docker образ собранный бинарный файл
COPY --from=0 /github/.bin/viewvrn .
COPY templates/logo.png templates/images/logo.png
COPY templates/viewpage.html templates/viewpage.html
COPY cron /etc/cron.d/viewvrn

COPY script.sh /script.sh
RUN chmod +x /script.sh 
#копируем каталог с картинками музея

#открываем порт для запросов к серверу TG
EXPOSE 8082
#запускаем приложение с nohup
#CMD ["./gettheatre"]
CMD ["/bin/bash", "-c", "/script.sh && date && chmod 644 /etc/cron.d/viewvrn && cron && tail -f /var/log/cron.log"]
#CMD service cron start && date && tail -f /var/log/syslog 