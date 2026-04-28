<div align="center">
  <img src="./frontend/public/favicon.svg" alt="interpolation_logo" width="100" height="100" />
  <br/>
  <br/>
  <div style="display: flex; justify-content: center; gap: 8px; flex-wrap: wrap;">
    <img alt="Go" src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" />
    <img alt="Gin" src="https://img.shields.io/badge/Gin-1.9+-00ADD8?logo=go&logoColor=white" />
    <img alt="Vue.js" src="https://img.shields.io/badge/Vue.js-3.4-4FC08D?logo=vuedotjs&logoColor=white" />
    <img alt="Vite" src="https://img.shields.io/badge/Vite-5-646CFF?logo=vite&logoColor=white" />
    <img alt="Vuetify" src="https://img.shields.io/badge/Vuetify-3-1867C0?logo=vuetify&logoColor=white" />
    <img alt="Chart.js" src="https://img.shields.io/badge/Chart.js-4-FF6384?logo=chartdotjs&logoColor=white" />
    <img alt="Docker" src="https://img.shields.io/badge/Docker-compose-2496ED?logo=docker&logoColor=white" />
    <img alt="Redis" src="https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white" />
    <img alt="RabbitMQ" src="https://img.shields.io/badge/RabbitMQ-3-FF6600?logo=rabbitmq&logoColor=white" />
  </div>
</div>

# Обучающее приложение "Численные методы. Интерполяция"

Интерактивное веб-приложение для изучения методов интерполяции (линейная, Лагранж, Ньютон).

## Архитектура

Проект построен на микросервисной архитектуре:

- **API Gateway**: Маршрутизация, CORS, Rate Limiting (Go + Gin).
- **Interpolation Service**: Вычислительное ядро с кэшированием (Go + Redis).
- **History Service**: Сохранение истории вычислений (Go + SQLite + RabbitMQ).
- **Frontend**: Интерактивный SPA (Vue 3 + Vuetify + Chart.js).

## Запуск через Docker Compose

```bash
docker compose up --build
```

## Доступные эндпоинты

- Фронтенд: `http://localhost/`
- API Gateway: `http://localhost/api/v1/`
- Swagger: `http://localhost/swagger/index.html`
- RabbitMQ Management: `http://localhost:15672` (guest/guest)

## Технологии

- **Backend**: Go 1.21, Gin, Redis, RabbitMQ, SQLite.
- **Frontend**: Vue 3, Vuetify 3, Chart.js.
- **DevOps**: Docker, Docker Compose, Nginx.

---

<div align="center">
  <img src="./frontend/public/favicon.svg" alt="interpolation_logo" width="100" height="100" />
</div>
