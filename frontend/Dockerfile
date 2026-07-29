FROM node:24-alpine AS dev

WORKDIR /app

COPY package.json package-lock.json* ./

RUN --mount=type=cache,target=/root/.npm \
    npm install

COPY . .

EXPOSE 3000

CMD ["npm", "run", "dev"]