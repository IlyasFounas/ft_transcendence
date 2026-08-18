#!/bin/sh
/app/bun_migrator init
/app/bun_migrator migrate
exec /app/golang_api