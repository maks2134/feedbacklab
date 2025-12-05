# Keycloak Authentication Guide

## Обзор

Приложение использует Keycloak для аутентификации и авторизации пользователей. Все API endpoints (кроме health check и swagger) требуют JWT токен в заголовке Authorization.

## Настройка

### Переменные окружения

Добавьте следующие переменные в `.env.docker`:

```env
KEYCLOAK_URL=http://localhost:8082
KEYCLOAK_REALM=feedbacklab
KEYCLOAK_CLIENT_ID=feedbacklab-api
KEYCLOAK_ADMIN=admin
KEYCLOAK_ADMIN_PASSWORD=admin
```

## Получение токена

### Для пользователя

```bash
curl -X POST http://localhost:8082/realms/feedbacklab/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=feedbacklab-api" \
  -d "username=admin" \
  -d "password=admin123"
```

Ответ:
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ...",
  "expires_in": 300,
  "refresh_expires_in": 1800,
  "token_type": "Bearer"
}
```

## API Endpoints

### 1. Создание пользователя (только для администратора)

**POST** `/api/auth/users`

Требует роль `admin`.

**Request:**
```json
{
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "password": "temporary123"
}
```

**Response (201):**
```json
{
  "userId": "uuid-here",
  "email": "user@example.com",
  "username": "user",
  "message": "User created successfully. Password is temporary and must be changed on first login."
}
```

**Пример:**
```bash
curl -X POST http://localhost:8080/api/auth/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "password": "temporary123"
  }'
```

### 2. Получение информации о текущем пользователе

**GET** `/api/auth/me`

Требует аутентификации.

**Response (200):**
```json
{
  "id": "uuid-here",
  "username": "user",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "emailVerified": false,
  "enabled": true
}
```

**Пример:**
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Смена пароля

**POST** `/api/auth/password/change`

Требует аутентификации.

**Request:**
```json
{
  "currentPassword": "oldpassword123",
  "newPassword": "newpassword456"
}
```

**Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

**Пример:**
```bash
curl -X POST http://localhost:8080/api/auth/password/change \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "currentPassword": "oldpassword123",
    "newPassword": "newpassword456"
  }'
```

## Роли

- **admin** - Администратор с полным доступом (может создавать пользователей)
- **user** - Обычный пользователь

## Использование токена

Все защищенные endpoints требуют заголовок:

```
Authorization: Bearer <access_token>
```

## Первый вход пользователя

1. Администратор создает пользователя через `/api/auth/users`
2. Пользователю выдается временный пароль
3. Пользователь получает токен с этим паролем
4. Пользователь должен сменить пароль через `/api/auth/password/change`
5. После смены пароля пользователь может использовать новый пароль для получения токенов

## Ошибки

### 401 Unauthorized
- Токен отсутствует или невалиден
- Токен истек

### 403 Forbidden
- У пользователя нет необходимой роли

### 400 Bad Request
- Неверный формат запроса
- Текущий пароль неверен (при смене пароля)

### 409 Conflict
- Пользователь с таким email уже существует

