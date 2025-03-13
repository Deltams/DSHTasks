# Шаги для развертывания приложения
## Требования
Перед началом работы убедитесь, что у вас установлены следующие инструменты:
- `Docker`: Для сборки и запуска контейнера приложения.
- `Git`: Для клонирования репозитория с проектом.

## Клонирование проекта
Скопируйте все файлы проекта в одну директорию. Это можно сделать с помощью команды `git clone`
\```
git clone https://github.com/Deltams/DSHTasks.git
\```

## Настройка конфигурации базы данных
Откройте файл `config/data_base/DBConnect.json` и отредактируйте параметры подключения к базе данных. Убедитесь, что указаны правильные данные для доступа к вашей базе данных.

## Запуск приложения
1. Откройте терминал и перейдите в корневую директорию проекта:
\```
cd DSHTasks
\```

2. Соберите и запустите контейнеры с помощью Docker Compose:
\```
docker-compose up --build
\```

3. Эта команда выполнит сборку образа и запустит приложение на порте 8080. Если вы хотите запустить приложение повторно, используйте команду:
\```
docker-compose up
\```

4. После успешного запуска приложения, оно будет доступно по адресу: `http://localhost:8080`.

## Доступ к панели администрирования базы данных (PGAdmin)
Панель администратора PGAdmin доступна по адресу: `http://localhost:4080`. Логин и пароль для входа указаны в файле `docker-compose.yml`.

## База данных
При первом запуске (`docker-compose up --build`) будет создана база данных, описанная в файле `init.sql`. По умолчанию она развернута для внешнего соединения по порту 4032.

## Использование API
### Примеры запросов к API:

### Создание новой задачи
\```
curl -X POST http://localhost:8080/tasks \
-H 'Content-Type: application/json' \
-d '{"title": "Новая задача", "description": "Описание задачи"}'
\```

### Получение списка всех задач
\```
curl -X GET http://localhost:8080/tasks
\```

### Получение конкретной задачи по ID
\```
curl -X GET http://localhost:8080/tasks/1
\```

### Обновление задачи
\```
curl -X PUT http://localhost:8080/tasks/1 \
-H 'Content-Type: application/json' \
-d '{"title": "Обновленная задача", "description": "Новое описание задачи"}'
\```

### Удаление задачи
\```
curl -X DELETE http://localhost:8080/tasks/1
\```
