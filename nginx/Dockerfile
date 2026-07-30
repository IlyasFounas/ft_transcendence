FROM nginx:stable-alpine AS dev
COPY conf/dev.conf /etc/nginx/conf.d/dev.conf
EXPOSE 5130

FROM nginx:stable-alpine AS prod
COPY conf/prod.conf /etc/nginx/conf.d/prod.conf
EXPOSE 80
EXPOSE 443