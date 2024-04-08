# authentication-deanery

TODO: 09.04.2024

* Необходимо добавить страницу для модерации (перенести туда регистрацию + реализовать добавление/удаление ролей через UI) + оформить главную страницу.
* Также сделать редирект на страницу логина, если у человека отсутсвует актуальный jwt токен.

## configuration file example

```
UserDatabase:
  user: 'postgres'
  password: 'password'
  host: 'localhost'
  port: '5432'
  dbname: 'authentication'
  sslmode: 'disable'
  maxconns: '30'
  
UserService:
  bindAddr: '0.0.0.0:8000'
  readTimeout: '1'
  writeTimeout: '1'
  idleTimeout: '5'

Encryption:
  passwordLength: '14'
  aesEncryptionKey: ''

Authentication:
  jwtSecretToken: 'secret'

EmailService:
  smtpEmail: ''
  smtpPassword: ''
  smtpService: ''
  smtpAddress: ''

Gist:
  url: ''
```
