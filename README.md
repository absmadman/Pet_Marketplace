# VK_Internship_Marketplace
REST API проект в рамках VK Internship 

# Быстрый старт
Для быстрого старта вам нужно иметь установленный Docker и docker-compose, перейти в корневую \
папку проекта(VK_Internship_Marketplace) и выполнить команду `make docker-compose-api`, проект развернется \
и будет принимать запросы по адресу `localhost:8080`

# Установка и запуск
Проект запускается в Docker контейнере с помощью docker-compose, для удобства Makefile содержит следующие зависимости:

`make help` - выводит подсказки по зависимостям \
`make docker-compose-api` - запускает приложение в контейнере с помощью compose \
`make clean-data` - очищает данные из базы данных и кеша\
`make docker-stop-api` - останавливает работу контейнеров\
`make docker-clean-api` - удаляет контейнеры\
`make server-logs` - выводит логи с контейнера сервера \
`make database-logs` - выводит логи с контейнера базы данных\
`make cache-logs` - выводит логи с контейнера для кеша\
`make all-logs` - выводит все логи вместе

# О проекте
Проект создан с помощью пакета Gin, запускается через -compose, помимо основного контейнера, проекту требуется "поднять" еще контейнер с postgresql и redis, это происходит автоматически через docker-compose. \
Соединение между контейнерами осуществляется с помощью links в docker-compose файле и переменных окружения

По мере разработки проекта я старался придерживаться чистой архитектуры, принципов SOLID, старался создать проект легко масштабируемым и поддерживаемым.\
Константы вынесены в файлы templates.go

Регулировка различных параметров происходит с помощью параметров в файле ***.env***.\

Проект реализует следующие эндпоинты и методы с примерами:
* /auth
  * /register
      * POST без параметров, принимает json с полями login(строка) и password(строка), возвращает данные пользователя в случае успешной регистрации:
          * Например:
          ```
          curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d {'"login": "Dmitriy123", "password": Password1'}
          ```
          * Вернет json следующего вида:
          ```
          {
            "id" : 2,
            "login": "Dmitriy123",
            "password": "Password1"
          }
          ```
  * /login
      * GET без параметров, принимает json с полями login(строка) и password(строка), возвращает jwt токен по id пользователя в случае успешной авторизации:
        * Например:
        ```
        curl -X GET http://localhost:8080/auth/login -H "Content-Type: application/json" -d {'"login": "Dmitriy123", "password": Password1'}
        ```
        * Вернет строку следующего вида:
        ```
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE1NjA3ODcsImlhdCI6MTcxMTU1NzE4NywiSWQiOjJ9.qaOjwUaglWkQL9eH1vHGBsJnSvjcjJflruhOMIpoESg"
        ```
* /api
  * /advert
    * GET с параметром advert_id(число), этот запрос работает для авторизованных и неавторизованных пользователей
        * Например: 
        ```
        curl -X GET 'http://localhost:8080/api/advert?advert_id=1 
        ```
        * В случае успеха вернет сообщение:
        ```
        {
          "id": 1,
          "user_id": 2,
          "header": "Test advert",
          "text": "Test advert text",
          "image_url": "test/test_image_url",
          "address": "Test address",
          "price": 1.5,
          "datetime": "2024-03-27T16:40:51.834281Z",
          "by_this_user": false
        } 
        ```
    * POST без параметров, принимает json с полями header, text, image_url, address, price, только для авторизованных пользователей
        * Например:
        ```
        curl -X POST http://localhost:8080/api/advert -H "Authorization: Bearer (токен_пользователя)" -H "Content-Type: application/json" -d
        {
          '"header": "Test advert",
          "text": "Test advert text",
          "image_url": "test/test_image_url",
          "address": "Test address",
          "price": 1.5'
        }
        ```
        * Вернет:
        ```
        {
          "id": 1,
          "user_id": 2,
          "header": "Test advert",
          "text": "Test advert text",
          "image_url": "test/test_image_url",
          "address": "Test address",
          "price": 1.5,
          "datetime": "0001-01-01T00:00:00Z",
          "by_this_user": true
        } 
        ```
    * PUT с параметром advert_id а также json с полями которые следует обновить только для авторизованных и если объявление принадлежит этому пользователю
        * Например:
        ```
        curl -X PUT 'http://localhost:8080/api/advert?advert_id=1' -H "Authorization: Bearer (токен_пользователя) -H "Content-Type: application/json" -d
        {
          "text": "Put test",
          "image_url": "put/test/image_url",
          "address": "Put test address",
          "price": 100
        }
        ```
        * Вернет:
        ```
        {
          "id": 1,
          "user_id": 1,
          "header": "Test advert",
          "text": "Put test",
          "image_url": "put/test/image_url",
          "address": "Put test address",
          "price": 100,
          "datetime": "2024-03-27T16:40:51.834281Z",
          "by_this_user": true
        }
        ```
    * DELETE с параметром advert_id только для авторизованных и если объявление принадлежит этому пользователю
        * Например:
        ```
        curl -X DELETE 'http://localhost:8080/api/advert?advert_id=1' -H "Authorization: Bearer (токен_пользователя)" 
        ```
        * Вернет:
        ```
        {
          "id": 1,
          "user_id": 1,
          "header": "Test advert",
          "text": "Put test",
          "image_url": "put/test/image_url",
          "address": "Put test address",
          "price": 100,
          "datetime": "2024-03-27T17:05:00.874232Z",
          "by_this_user": false
        } 
        ```
  * /feed - поддерживает запросы для авторизированных пользователей и нет, также поддерживает несколько видов фильтрации:
    * GET c json содержащим параметры фильтрации
      * Например (сортировка по цене в порядке возрастания):
      ```
      {
        "min_price": 1,
        "max_price": 100000,
        "by_price": true,
        "ascending_direction": true
      } 
      ``` 
      * Вернет:
      ```
      {
      "feed": [
        {
          "id": 1,
          "user_id": 1,
          "header": "Test advert",
          "text": "Test advert text",
          "image_url": "test/test_image_url",
          "address": "Test address",
          "price": 1,
          "datetime": "2024-03-27T17:17:34.056469Z",
          "by_this_user": true
        },
        {
          "id": 2,
          "user_id": 1,
          "header": "Test advert",
          "text": "Test advert text",
          "image_url": "test/test_image_url",
          "address": "Test address",
          "price": 2,
        "datetime": "2024-03-27T17:17:36.562401Z",
        "by_this_user": true
        },
        {
          "id": 3,
          "user_id": 1,
          "header": "Test advert",
          "text": "Test advert text",
          "image_url": "test/test_image_url",
          "address": "Test address",
          "price": 3,
          "datetime": "2024-03-27T17:17:38.477777Z",
          "by_this_user": true
        },
        {
          "id": 4,
          "user_id": 1,
          "header": "Test advert",
          "text": "Test advert text",
          "image_url": "test/test_image_url",
          "address": "Test address",
          "price": 4,
          "datetime": "2024-03-27T17:17:40.457936Z",
          "by_this_user": true
        },
      ...
      ]
      }
      ```
    * Например (сортировка по времени добавления в порядке убывания):
      ```
      {
        "min_price": 1,
        "max_price": 100000,
        "by_price": false,
        "ascending_direction": false
      }  
      ``` 

# Функционал
* Проект реализован с помощью пакета Gin, который поддерживает много встроенного функционала
* Проект запускается в docker контейнере, всё необходимое можно конфигурировать с помощью файла ***.env***
* Проект структурно разбит на слои, что в будущем позволяет "безболезненно" и быстро заменять части программы, например базу данных и кеш
* Аутентификация пользователя происходит последством присвоения пользователю jwt токена сформированного из его id в базе данных
* Валидация:
  * Разрешены символы английского алфавита, цифры, а также - "!", "_", " ", "/"
  * Логин и пароль пользователя должен содержать минимум 8 символов, 1 заглавную букву, прописную букву и цифру, максимальная длина 36 символов
  * Текст обьявления должен содержать минимум 8, максимум 512 символов
  * Адрес, заголовок, ссылка на изображение должны содержать минимум 8, максимум 64 символа
  * Цена объявления - минимум 1, максимум 1 000 000
* Фильтрация ленты объявлений с помощью json структуры:
    ```
    {
      "min_price": минимальная цена объявления
      "max_price": максимальная цена обьявления,
      "by_price": сортировка по цене, если этот параметр равен false то сортировка происходит по дате,
      "ascending_direction": направление сортировки, если этот параметр равен false то сортировка происходит в убывающем порядке
    }
    ```
# Дополнительно
* Помимо основного задания, мной был реализован дополнитеьный функционал:
  *  Статическое конфигурирование через файл ***.env***
  *  Запуск в docker контейнере
  *  Кеширование данных с помощью отдельного контейнера с БД Redis, такой подход к кешированию позволит запускать проект в нескольких экземплярах с одним кешем
  *  Дополнительные методы у эндпоинта advert, такие как PUT и DELETE