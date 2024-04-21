Test Case 1: Registration with Missing Login
```
curl --location 'http://localhost:5441/register' \
--header 'Content-Type: application/json' \
--data '{"password": "passwordwithoutlogin"}'
```

Test Case 2: Registration with Missing Password
```
curl --location 'http://localhost:5441/register' \
--header 'Content-Type: application/json' \
--data '{"login": "loginwithoutpassword"}'

```

Test Case 3: Registration with Empty Body
```
curl --location 'http://localhost:5441/register' \
--header 'Content-Type: application/json' \
--data '{}'
```

Test Case 4: Registration with Special Characters in Login
```
curl --location 'http://localhost:5441/register' \
--header 'Content-Type: application/json' \
--data '{"login": "user!@#","password": "passwordForUser"}'
```