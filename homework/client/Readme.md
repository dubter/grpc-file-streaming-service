# Files Transfer Client

### Конфигурация:
Пример:
```
# Таймаут запроса клиента
timeout: 5s
# Порт gRPC сервера, к которому подключаемся
port: 50051

# Доступные методы gRPC сервера
get_file_list:
  exec: true

get_file:
  # Хотим ли запускать метод
  exec: true
  # Параметры запуска
  name: test.txt

get_file_info:
  exec: true
  name: test.txt
```

### Алгоритм запуска:
* *make build* 
* *cd build*
* *./file_transfer_client --config_path=../config/client.yaml* (либо путь к своему конфигу)
