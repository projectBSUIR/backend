# BetterSolve API

  ## Скачивание последней версии Go Fiber

  [Здесь](https://go.dev/doc/install) нужно скачать последнюю версию Go
  
  Далее необходимо склонировать репозиторий:
  ```
  git clone https://github.com/projectBSUIR/backend
  git checkout dev
  ```
  
  [Здесь](https://learn.microsoft.com/ru-ru/azure/developer/go/configure-visual-studio-code) можно установить Go для работы с ним в Visual Studio Code
  
  Далее необходимо установить в папке проекта Go Fiber
  ```
  go get github.com/gofiber/fiber/v2
  ```
  
  ## Запуск проекта
  
  Следующая команда запускает сервер на http://127.0.0.1:5000/
  
  ```
  go run main.go
  ```
# Endpoints

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
   
   Поля login, password и email имеют тип string
