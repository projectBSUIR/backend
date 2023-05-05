# BetterSolve API

  ## Скачивание последней версии Go Fiber

  [Здесь](https://go.dev/doc/install) нужно скачать последнюю версию Go
  
  Далее необходимо склонировать репозиторий:
  ```
  git clone https://github.com/projectBSUIR/backend
  git checkout main
  ```
  
  [Здесь](https://learn.microsoft.com/ru-ru/azure/developer/go/configure-visual-studio-code) можно установить Go для работы с ним в Visual Studio Code
  
  Далее необходимо установить в папке проекта необходимые модули 
  ```
  go get github.com/gofiber/fiber/v2
  go get github.com/go-sql-driver/mysql
  go get github.com/golang-jwt/jwt/v4
  ```
  ## Скачивание MySQL
  
  [Здесь](https://dev.mysql.com/downloads/installer/) нужно скачать установщик для mysql.
  
  При его установке выбирать все дефолтные конфиги.
  
  Логин и пароль должны быть `root` и `password` соответственно.
  
  ### Создание Базы Данных
  
  Необходимо открыть MySQL Workbench. Создать новую схему с любым названием. Далее вставить в SQL скрипт текст из discord-сервера в канале `бд` и его запустить, нажав на молнию. 
  
  ### Запуск сервера MySQL
  
  Необходимо запустить MySQL Installer и в строке MySQL сервер выбрать Reconfigure. Нажимать далее до момента ввода пароля. После его ввода нажать на check и нажать далее. На следующей странице нажать execute и ждать запуска сервера. 

  ## Запуск проекта
  
  Следующая команда запускает сервер на http://127.0.0.1:5000/
  
  ```
  go run main.go
  ```

  ## Регистрация и авторизация пользователя
   ### Регистрация [POST]
    
   По следующей ссылке можно зарегистрировать пользователя в случае, если не существует пользователя с таким же полем login
   
   ```
   http://127.0.0.1:5000/register
   ```
   
   Для успешной регистрации пользователя необходимо передавать в запрос JSON следующего вида:
   
   ```
   {
        "login": login,
        "password": password,
        "email": email
   }
   ```
   
   Регистрация прошла успешно, если статус запроса 200
   
   В любых других случаях в Response Body будет возвращён JSON:
   ```
   {
         "message": message
   }
   ```
   
   Типы:
   
   message, login, password, email - string
  
   ### Авторизация [POST]
  
   По следующей ссылке можно авторизовать пользователя в случае, если он существует и введён верно пароль:
   ```
   http://127.0.0.1:5000/login
   ```
   
   Для успешной авторизации пользователя необходимо передавать в запрос JSON следующего вида:
   
   ```
   {
        "login": login,
        "password": password,
   }
   ```
   
   После успешной авторизации(статус 200) в Response Body будет возвращён JSON следующего вида:
   
   ```
   {
         "access_token": access_token
   }
   ```
   
   В любых других случаях в Response Body будет возвращён JSON:
   ```
   {
         "message": message
   }
   ```
   
   Типы:
   
   access_token, message, login, password - string
   
   ### Logout [POST]
   
   По следующей ссылке можно выйти из аккаунта в случае, если вы были авторизованы
   ```
   http://127.0.0.1:5000/logout
   ```
   
   Выход является успешным, если статус запроса 200.
   
   В противных случаях. в Response Body будет возвращён JSON:
   ```
   {
         "message": message
   }
   ```
   
   Типы:
   
   message - string
   
   
