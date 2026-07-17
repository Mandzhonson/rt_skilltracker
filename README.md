# SkillTracker

SkillTracker — это веб-приложение для управления развитием профессиональных навыков сотрудников.

Система позволяет:

- создавать планы развития сотрудников;
- автоматически генерировать планы обучения с помощью ИИ (Ollama);
- выполнять теоретические тесты;
- проверять практические задания;
- отслеживать прогресс сотрудников;
- управлять пользователями, ролями и менеджерами;
- хранить аватары пользователей в MinIO.

---

# Технологии

Backend:

- Go
- Gin
- PostgreSQL
- Redis
- MinIO
- JWT
- Ollama

Frontend:

- React
- Vite

Контейнеризация:

- Docker
- Docker Compose

---

# Запуск проекта

## 1. Клонировать репозиторий

```bash
git clone https://github.com/Mandzhonson/rt_skilltracker.git
cd skilltracker
```

---

## 2. Установить Ollama

Скачать можно с официального сайта:

https://ollama.com/download

---

## 3. Скачать модель

Например:

```bash
ollama pull qwen3:14b
```

---

## 4. Запустить Ollama

Backend работает с Ollama, запущенной на хостовой машине.

Запуск:

```bash
OLLAMA_HOST=0.0.0.0:11434 ollama serve
```

Проверить работу можно:

```bash
curl http://localhost:11434/api/tags
```

---

## 5. Создать файл `.env`

Пример:

```env
HTTP_PORT=8080
HTTP_HOST=0.0.0.0
HTTP_READ_TIMEOUT=5s
HTTP_WRITE_TIMEOUT=300s
HTTP_SHUTDOWN_TIMEOUT=10s

POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=skilltracker
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

ADMIN_EMAIL=admin@skilltracker.local
ADMIN_PASSWORD=admin123456
ADMIN_FIRST_NAME=System
ADMIN_LAST_NAME=Administrator

REDIS_PASSWORD=redispass
REDIS_HOST=redis
REDIS_PORT=6379

JWT_ACCESS_SECRET=super-secret-access-key
JWT_REFRESH_SECRET=super-secret-refresh-key
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

OLLAMA_URL=http://host.docker.internal:11434
OLLAMA_MODEL=qwen3:14b
OLLAMA_TIMEOUT=120s

MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_HOST=minio
MINIO_API_PORT=9000
MINIO_BUCKET=avatar
```

---

## 6. Запустить проект

```bash
make up
```

Будут автоматически:

- запущены PostgreSQL;
- запущен Redis;
- запущен MinIO;
- применены миграции базы данных;
- собран и запущен backend;
- собран и запущен frontend.

---

# Makefile

Запуск проекта

```bash
make up
```

Остановка

```bash
make down
```

Удаление контейнеров и томов

```bash
make clear
```

Применить миграции вручную

```bash
make migrate-up
```

Откатить миграции

```bash
make migrate-down
```

Посмотреть текущую версию миграций

```bash
make migrate-version
```

---

# Доступные сервисы

Backend

```
http://localhost:8080
```
Swagger API

```
http://localhost:8080/swagger/index.html
```

Frontend

```
http://localhost:3000
```

MinIO Console

```
http://localhost:9001
```

---

# Учетные данные администратора

После первого запуска автоматически создается администратор.

```
Email:
admin@skilltracker.local

Password:
admin123456
```

---

# Примечания

- Ollama должна быть запущена **до старта Docker Compose**.
- Backend обращается к Ollama через

```
http://host.docker.internal:11434
```
