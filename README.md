# План
## Создание токенов:
2. Создать пару
1. Проверить, refresh токен юзера в базе
3. Сохранить (обновить, если уже есть) refresh токен в базе (хэш)
4. Отдать результат

## Обновление токенов:
1. Принять refresh токен и user_id
2. Сходить в базу и найти токен для user_id
3. Проверить сходство (из базы хэш)
4. Создать пару
5. Обновить refresh токен в базе по user_id (хэш)
5. Отдать результат