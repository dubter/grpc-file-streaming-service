# Files Transfer Server

### Архитектура:
* HTTP Gateway - принимает запросы и проксирует в GRPC Server
* GRPC Server - валидирует и обрабытывает запросы, храня файлы в локальной файловой системе
* Чтобы посмотреть Swagger, перейдите в `/docs/file_transfer.swagger.json` и вставьте json сюда [swagger editor](https://editor.swagger.io/)

### Конфигурация:
Пример:
```
# Директория, в которой GRPC Сервер будет хранить файлы
folder_path: ~/data

http_gateway:
    port: 8080
    read_timeout: 30s
    write_timeout: 30s
    # Таймаут, в течение которого сервер ждёт пока все соединения закроются:
    shutdown_timeout: 30s

grpc_server:
    port: 50051
    # Размер одного блока данных, передаваемого по сети в байтах. В данном случае 64 * 1024 = 64 Кб:
    chunk_size: 65536
```

### Алгоритм запуска:
* *make build* 
* *cd build*
* *./file_transfer_server --config_path=../config/server.yaml* (либо путь к своему конфигу)
