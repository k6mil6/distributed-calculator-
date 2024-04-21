Test Case 1: Authorization with Incorrect Password
```
curl --location "http://localhost:5441/login" ^
--header "Content-Type: application/json" ^
--data "{\"login\": \"test_user\",\"password\": \"incorrect_password\"}"
```

Test Case 2: Authorization with Missing Login
```
curl --location "http://localhost:5441/login" ^
--header "Content-Type: application/json" ^
--data "{\"password\": \"passwordwithoutlogin\"}"
```

Test Case 3: Authorization with Empty Body
```
curl --location "http://localhost:5441/login" ^
--header "Content-Type: application/json" ^
--data "{}"
```

Test Case 4: Authorization with Special Characters in Password
```
curl --location "http://localhost:5441/login" ^
--header "Content-Type: application/json" ^
--data "{\"login\": \"userWithSpecialChars\",\"password\": \"password!@#\"}"
```
