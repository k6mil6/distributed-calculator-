Test Case 1: Empty Expression
```
curl --location "http://localhost:5441/calculate" ^
--header "Content-Type: application/json" ^
--header "Authorization: Bearer test_token" ^
--data "{\"id\": \"1a2b3c4d-1234-5678-90ab-cdef12345678\", \"expression\": \"\"}"
```

Test Case 2: Calculation Resulting in a Large Number
```
curl --location "http://localhost:5441/calculate" ^
--header "Content-Type: application/json" ^
--header "Authorization: Bearer test_token" ^
--data "{\"id\": \"1a2b3c4d-1234-5678-90ab-cdef12345678\", \"expression\": \"10000*10000\"}"
```
