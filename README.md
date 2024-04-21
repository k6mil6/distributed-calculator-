# По всем вопросам tg - @kamil_66

[Диаграмма](https://excalidraw.com/#json=iTZfpCq_xgxifx-g8BRbR,VzcX5ZcA6IDuZUEc26rx5A), как тут всё работает

Компоненты:
- migrator - для применения миграций
- http_server - для взаимодействия с пользователем
- orchestrator - сохраняет, делит выражения на подвыражения, предоставляет воркерам подвыражения по grpc
- agent - содержит в себе указанное в конфиге кол-во воркеров 

Для проверки задания потребуется установленные 

- [docker](https://docs.docker.com/get-docker/) 
- [git](https://git-scm.com/downloads)
- [postman](https://www.postman.com/downloads/) - для удобства отправки запросов, по желанию 
(в папке examples/postman лежит файл, который можно импортировать в постман и получить примеры запросов, в postman jwt token устанавливается при настройке запроса во вкладке Authorisation, type bearer token)

Для развертывания приложения необходимо:

1. клонировать репозиторий в удобную папку (git clone https://github.com/k6mil6/distributed-calculator.git .)
2. перейти в папку
3. в консоле прописать docker-compose up (если оркестратор не запустился, необходимо перезапустить его, либо через интерфейс docker desktop, либо нажать ctrl+c и написать docker-compose еще раз)
4. ниже представлены примеры для проверки работоспособности (в папке examples/curl лежат последующие примеры, но с разными значениями)

пример регистрации пользователя (mac/linux)
```
curl --location 'http://localhost:5441/register' \
--header 'Content-Type: application/json' \
--data '{"login": "your_login","password": "your_password"}'
```

windows (cmd)
```
curl --location "http://localhost:5441/register" ^
--header "Content-Type: application/json" ^
--data "{\"login\": \"your_login\",\"password\": \"your_password\"}"
```

пример авторизации (mac/linux)
```
curl --location 'http://localhost:5441/login' \
--header 'Content-Type: application/json' \
--data '{"login": "your_login","password": "your_password"}'
```

windows (cmd)
```
curl --location "http://localhost:5441/login" ^
--header "Content-Type: application/json" ^
--data "{\"login\": \"kamil\",\"password\": \"12345\"}"
```

пример вычисления выражения с передачей таймаутов (mac/linux)
(token - получаем в ответе на авторизацию)
```
curl --location 'http://localhost:5441/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your_token' \
--data '{"id": "3422b448-2460-4fd2-9183-8000de6f8348","expression": "2+2+3","timeouts": {"+": 10,"-": 20,"/": 10,"*": 5}}'
```

windows (cmd)
```
curl --location "http://localhost:5441/calculate" ^
--header "Content-Type: application/json" ^
--header "Authorization: Bearer your_token" ^
--data "{\"id\": \"e58ed763-928c-4155-bee9-fdbaaadc15f4\", \"expression\": \"2-2*2*10000\",\"timeouts\": {\"+\": 10,\"-\": 20,\"/\": 20,\"*\": 20}}"
```

timeouts - параметр, который можно не указывать, будет использовано последнее добавленное значение
также важно, чтобы id был формата uuid, для проверки нескольких выражений можно изменять последние цифры самого id

пример вычисления выражения без передачи таймаутов(mac/linux)
```
curl --location 'http://localhost:5441/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your_token' \
--data '{"id": "3422b448-2460-4fd2-9183-8000de6f8348","expression": "2+2+3"}'
```
windows(cmd)
```
curl --location "http://localhost:5441/calculate" ^
--header "Content-Type: application/json" ^
--header "Authorization: Bearer your_token" ^
--data "{\"id\": \"e58ed763-928c-4155-bee9-fdbaaadc15f4\", \"expression\": \"2-2*2*10000\"}"
```

пример получения выражения по id (macos/linux)

```
curl --location 'http://localhost:5441/expression/3422b448-2460-4fd2-9183-8000de6f8346' \
--header 'Authorization: Bearer your_token'
```

windows(powershell)
```
curl --location "http://localhost:5441/expression/3422b448-2460-4fd2-9183-8000de6f8346" ^
--header "Authorization: Bearer your_token"
```

пример получения всех выражений (macos/linux)

```
curl --location 'http://localhost:5441/all_expressions' \
--header 'Authorization: Bearer your_token'
```

windows(cmd)
```
curl --location "http://localhost:5441/all_expressions" ^
--header "Authorization: Bearer your_token"
```

пример установки таймаутов для операций (macos/linux)

```
curl --location 'http://localhost:5441/set_timeouts' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer your_token' \
--data '{"timeouts": {"+": 10, "-": 10,"*": 10,"/": 10}}'
```

windows(cmd)
```
curl --location "http://localhost:5441/set_timeouts" ^
--header "Content-Type: application/json" ^
--header "Authorization: Bearer your_token" ^
--data "{\"timeouts\": {\"+\": 10,\"-\": 10,\"*\": 10,\"/\": 10}}"
```

пример получения актуальных таймаутов (macos/linux)
```
curl --location 'http://localhost:5441/actual_timeouts' \
--header 'Authorization: Bearer your_token' 
```

windows(cmd)
```
curl --location "http://localhost:5441/actual_timeouts" ^
--header "Authorization: Bearer your_token" 
```




