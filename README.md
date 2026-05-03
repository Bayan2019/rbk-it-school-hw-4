# Домашнее задание 4

## Database Schema

[![rbk](./rbk-weather-api.svg "rbk DataBase")](https://dbdiagram.io/d/rbk-weather-api-69ef55e7c6a36f9c1b925ccc)

## REST API

### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"bayan@example.com","password":"tramp",
        "first_name": "Bayan", "last_name": "User",
        "is_active": true}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"bayan@example.com","password":"tramp"}'
```

### Profile

```bash
curl http://localhost:8080/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Добавить город пользователю

```bash
curl -X POST http://localhost:8080/cities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "city": "Almaty"
  }'
```

### Список городов пользователя

```bash
curl -X GET http://localhost:8080/cities \
-H "Authorization: Bearer $TOKEN"
```

### Удалить город

```bash
curl -X DELETE http://localhost:8080/cities/3 \
-H "Authorization: Bearer $TOKEN"
```

### Получить погоду по всем городам пользователя

```bash
curl -X GET http://localhost:8080/weather \
-H "Authorization: Bearer $TOKEN"
```

### Получить историю с фильтрацией

```bash
curl -X GET http://localhost:8080/weather/history?city=astana \ 
-H "Authorization: Bearer $TOKEN"
```



### Admin endpoint

```bash
curl -X GET http://localhost:8080/users \
-H "Authorization: Bearer $TOKEN"
```

```bash
curl -X GET http://localhost:8080/users/3 \
-H "Authorization: Bearer $TOKEN"
```

```bash
curl -X DELETE http://localhost:8080/users \
-H "Authorization: Bearer $TOKEN"
```


## References 

- [Database Schema](https://dbdiagram.io/d/rbk-weather-api-69ef55e7c6a36f9c1b925ccc)
- [Json to Go Struct](https://transform.tools/json-to-go)
- [Geocoding API](https://open-meteo.com/en/docs/geocoding-api)
- [Welcome to Nominatim](https://nominatim.openstreetmap.org/ui/search.html)
- [Bcrypt Hash Generator](https://bcrypt-generator.com)