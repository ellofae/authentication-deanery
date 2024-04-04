# authentication-deanery

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
  aesEncryptionKey: 'supersecretaesencryptionkeyhashd'
```