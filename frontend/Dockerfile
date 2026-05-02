FROM node:20-alpine AS builder

WORKDIR /app
ENV npm_config_loglevel=warn
ENV npm_config_fund=false
ENV npm_config_audit=false
ENV npm_config_update_notifier=false
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
# Set build-time env if needed
ARG VITE_API_URL=/api/v1
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build -- --logLevel warn

FROM nginx:alpine
COPY docker/nginx-main.conf /etc/nginx/nginx.conf
COPY --from=builder /app/dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
