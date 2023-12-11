
## Подготовка к запуску

1. Создайте файл `.env` в папке `/config` с необходимыми данными. Пример содержания файла:

    ```env
   POSTGRES_HOST=postgres
   POSTGRES_PORT=5432
   POSTGRES_USER=test
   POSTGRES_PASSWORD=1234
   POSTGRES_DB=urls_db
   POSTGRES_SSLMODE=disable
   
   REDIS_HOST=redis
   REDIS_PORT=6379
   REDIS_PASSWORD=1234
   
   HTTP_SERVER_HOST=0.0.0.0
   HTTP_SERVER_PORT=8080
   
   GRPC_SERVER_HOST=0.0.0.0
   GRPC_SERVER_PORT=50051
    ```

2. Используемую БД можно выбрать в файле `config.yaml` в папке `/config` в поле `db_type`. 
Доступные БД: `postgres`, `redis`

## Запуск

### Запуск тестов и приложения
```bash
make full-run
```
### Запуск приложения
```bash
make run
```
### Запуск тестов
```bash
make test
```

## Алгоритм получения случайной строки
Генерирует случайную строку заданной длины. Использует пакет `crypto/rand` для генерации случайных чисел.
Для каждого символа выбирается случайный индекс из `alphabet`, и этот символ добавляется к результату.




