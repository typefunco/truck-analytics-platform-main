FROM nginx:latest

# Копируем конфигурацию Nginx
COPY nginx.conf /etc/nginx/nginx.conf

# Копируем файлы фронтенда в контейнер
COPY ./frontend /usr/share/nginx/html

# Открываем порт 80 для веб-сервера
EXPOSE 80