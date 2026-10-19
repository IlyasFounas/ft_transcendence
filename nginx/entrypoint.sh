#!/bin/sh

mkdir -p /etc/nginx/ssl
if [ ! -f /etc/nginx/ssl/transcendence.key ]; then
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/transcendence.key \
  -out /etc/nginx/ssl/transcendence.crt \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,DNS:campus.ocefunding.com,DNS:ocefunding.com,IP:127.0.0.1"
fi

exec nginx -g 'daemon off;'