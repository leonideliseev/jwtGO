# JWT APP

## Описание Проекта

Проект реализует логику создания и обновления JWT токенов для пользователя.
Refresh токен хранится в базе в виде bcrypt хэша.
В случае, если при обновлении токенов меняется ip, это проверяется ~~и пользователю отправляется уведомление~~:

```go
// internal/handler/tokens.go
if nowIP != wasIP {
	// логика отправки сообщения пользователю о смене ip
}
```

## Запуск

### Из `main`:

Для запуска приложения локально, используя `main.go` можно выполнить команду:

```
go run cmd/main.go
```

Для этого понадобится запущенная локально база данных `PostgreSQL`.

### Docker-compose:

В приложении написан `docker-compose.yml` для запуска приложения в докер-контейнере с образом базы.

```
docker-compose up -d --build
```

## настройка окружения

Приложение берёт переменные для запука из окружения проекта.
Доступные переменные окружения:

|`config.yaml`|`.env`|Default value|Type value|Description|
|:-:|:-:|:-:|:-:|---|
|`http.port`|`HTTP_PORT`|`8080`|`port`|Port on which the HTTP server will run.|
|`http.host`|`HTTP_HOST`|`0.0.0.0`|`host`|Host address for the HTTP server.|
|`postgresql.user`|`PG_USER`|`postgres`|`username`|Username for connecting to the PostgreSQL database.|
|`postgresql.password`|`PG_PASSWORD`|`password`|`password`|Password for connecting to the PostgreSQL database.|
|`postgresql.host`|`PG_HOST`|`0.0.0.0`|`host`|Host address of the PostgreSQL database.|
|`postgresql.port`|`PG_PORT`|`5432`|`port`|Port on which the PostgreSQL database is running.|
|`postgresql.database`|`PG_DATABASE`|`postgres`|`database`|Name of the PostgreSQL database.|
|`postgresql.ssl_mode`|`PG_SSL`|`disable`|`ssl mode`|SSL mode for the PostgreSQL connection (e.g., `disable`, `require`).|
|`jwt.access_sign_key`|`ACCESS_JWT_KEY`|`access_secret`|`key`|Secret key used to sign access JWT tokens.|
|`jwt.access_token_ttl`|`ACCESS_JWT_TTL`|`60m`|`duration`|Time-to-live for access JWT tokens.|
|`jwt.refresh_sign_key`|`REFRESH_JWT_KEY`|`refresh_secret`|`key`|Secret key used to sign refresh JWT tokens.|
|`jwt.refresh_token_ttl`|`REFRESH_JWT_TTL`|`168h`|`duration`|Time-to-live for refresh JWT tokens.|

Эти переменные можно задать в config.yaml/создать .env/прокинуть в контейнеры при запуске/оставить по умолчанию.
Приоритет следующий (от наибольшего к меньшему):

- В командной строке при запуске
- Из окружения (секция `environment` -> `.env` -> `export` в linux)
- `config.yaml`
- значения по умолчанию

За основу рекомендуется взять имеющийся `config.yaml` и переименовать `.env.example` в `.env` с заменой на необходимые значения.

**ATTENTION 1**: если запускать приложение _локально_, то `postgresql.host` (`PG_HOST`) должно быть `0.0.0.0`, но если через _docker-compose_, то `postgres`.

**ATTENTION 2**: Поля `PG_PASSWORD`, `PG_USER`, `PG_DATABASE`, `APP_PORT` вынесены в `docker-compose.yml` как ссылки `${}`, поэтому если убрать их из `.env`, то надо либо передавать их другим способ, либо указать конкретное значение.

## Используемые технологии и библиотеки

### Технологии

- Go
- PostgreSQL
- JWT
- Git
- Docker

### Библиотеки

- `github.com/gin-gonic/gin` - маршрутизация запросов на API
- `github.com/jackc/pgx/v5` - создание полключения к базе
- `github.com/Masterminds/squirrel` - создание sql запросов
- `github.com/golang-migrate/migrate/v4` - миграции
- `github.com/ilyakaznacheev/cleanenv` - чтение переменных окружения
- `github.com/golang-jwt/jwt` - работа с jwt
- `"github.com/stretchr/testify/assert"` - тестирование

### Сопутствующее

При разработке проекта старался учесть и применить подходы чистого кода, такие как:
`Dependency Injection`, `SOLID`, `KISS`, `DRY`

## API DOCS

### `/auth/tokens` [GET]

#### Request

`user_id`, заданный в параметрах url.

Тип UUID (GUID): `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`.

Пример: `auth/tokens?user_id=aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa`

#### Reaponse

`access_token` и `refresh_token`

#### Опсиание

По этому путь происходит создание `access` и `refresh` токенов для пользователя с `user_id`.

### `/auth/refresh` [GET]

#### Request

`user_id`, в параметрах url.

Тип UUID (GUID): `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`.

`refresh_token`, в параметрах url в кодировке base64.

Пример: `auth/tokens?user_id=aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa&refresh_token=very-big-refresh-token`

#### Reaponse

`access_token` и `refresh_token`

#### Опсиание

По этому путь происходит создание новых `access` и `refresh` токенов, по действующему `refresh` токену, для пользователя с `user_id`.
После этого использованный `refresh` токен становится недействительным.

Обновление токенов можно реализовать через `POST` метод.
Но в этом, как на мой взгляд, нет необходимости, так как если токен валидный, то он обновиться и больше будет не пригоден для использования.
Ну а если не валидный, то никаких проблем и так нет.

> Обычно, приложения с JWT токенами должны поддерживать _логику деактивации токенов_.
> Её можно реализовать, если хранить где-то валидные tokenID в связке с user_id.
> Тогда если пользователь захочет завершить все сеансы и разлогиниться, то можно удалить все хранимые tokenID у пользователя (можно просто добавить в базу ещё и поле user_id :)

## Тесты

Тесты пока в процессе написания, на данный момент есть только один.

Запуск тестов:

```
go test -v ./...
```
