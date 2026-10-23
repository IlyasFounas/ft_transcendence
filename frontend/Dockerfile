FROM node:24-alpine AS base
WORKDIR /app
COPY package.json package-lock.json* ./
COPY . .

FROM base AS dev
RUN --mount=type=cache,target=/root/.npm \
    npm install
EXPOSE 3000
CMD ["npm", "run", "dev"]

FROM base AS build

ARG VITE_MAPILLARY_ACCESS_TOKEN
ENV VITE_MAPILLARY_ACCESS_TOKEN=$VITE_MAPILLARY_ACCESS_TOKEN

RUN npm install
RUN npm run build

FROM nginx:stable-alpine AS prod
COPY --from=build /app/dist /usr/share/nginx/html
COPY ./front.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]