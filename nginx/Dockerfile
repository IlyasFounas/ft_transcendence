FROM nginx:stable-alpine AS dev
RUN rm -f /etc/nginx/conf.d/default.conf
COPY conf/dev.conf /etc/nginx/conf.d/dev.conf
EXPOSE 5130

FROM nginx:stable-alpine AS prod
WORKDIR /app
# Install app dependencies
RUN apk add --no-cache openssl
RUN mkdir -p /etc/nginx/ssl
RUN rm -f /etc/nginx/conf.d/default.conf
COPY conf/prod.conf /etc/nginx/conf.d/prod.conf
COPY entrypoint.sh ./
EXPOSE 443
EXPOSE 80
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]