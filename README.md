# Практическая работа №2 — Docker и Docker Compose

Учебное fullstack-приложение: Go API + статичный HTML/CSS/JS + PostgreSQL.

## Запуск

Нужен установленный Docker Desktop (Windows/macOS) или Docker Engine с Docker Compose v2.

1. Распакуйте архив в отдельную папку.
2. Откройте терминал в папке, где лежит `docker-compose.yml`.
3. Выполните:

   ```bash
   docker compose config -q
   docker compose up -d --build
   docker compose ps
   ```

4. Откройте http://localhost:8080
5. Проверьте API:

   ```bash
   curl http://localhost:8080/api/items
   ```

   Ожидается JSON-массив с двумя демонстрационными записями.

На Windows PowerShell, если `curl` работает как псевдоним, можно выполнить:
`Invoke-RestMethod http://localhost:8080/api/items`.

## Полезные команды

```bash
docker compose logs -f app
docker compose logs -f db
docker compose ps
docker compose down
```

`docker compose down` сохраняет именованный том с данными БД. Чтобы удалить и данные:
`docker compose down -v` (осторожно: данные будут удалены).

## Что показать на защите

1. `docker compose up -d --build` — сборка и запуск обоих сервисов.
2. `docker compose ps` — состояние контейнеров.
3. Открыть `http://localhost:8080` и показать список записей.
4. Выполнить запрос `curl http://localhost:8080/api/items`.
5. `docker compose logs app` — показать подключение приложения.
6. `docker image ls` — показать собранный образ приложения.

## Как выполнены критерии

- **Multi-stage:** `backend/Dockerfile` содержит этап `builder` и минимальный runtime `distroless`; в финальный образ копируется только бинарник.
- **Секреты:** имя пользователя и пароль находятся в `.env`, а не захардкожены в Go-коде или Compose. `.env` включён в `.gitignore`. Значения в архиве предназначены только для локальной учебной демонстрации.
- **PostgreSQL:** база данных работает отдельным сервисом `db`; SQLite не используется.
- **Два тома:** `db_data` — именованный том PostgreSQL; `./frontend:/app/static:ro` — bind mount папки хоста.
- **Порты:** наружу проброшен только веб-порт `APP_PORT` (по умолчанию 8080). У базы данных нет `ports`.
- **Compose:** у сервиса `app` задан `build`; поле `image` намеренно отсутствует.
- **Готовность БД:** `healthcheck` + `depends_on: condition: service_healthy`, а Go-приложение дополнительно повторяет подключение.
- **Стартовая таблица:** приложение создаёт таблицу и добавляет демонстрационные записи, если таблица пуста; отдельный init-volume не нужен.

## Если порт 8080 занят

Измените `APP_PORT` в `.env`, например на `8081`, затем откройте `http://localhost:8081`.

## Важно

Для сдачи на своём устройстве сохраните всю папку проекта. Docker должен быть установлен и запущен. При первом запуске Docker скачает базовые образы и Go-зависимость, поэтому нужен интернет.
