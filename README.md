Weather API
=============

This is a simple project just for fun

Requirements
-----------
go 1.14+ 


docker

Installation
-----------

```
docker pull redis:alpine
docker run --name weather-api-redis -p 6379:6379 -d redis:alpine
```

and from root directory: 
```
go run main.go
```


Usage
-----
You can find a current weather for a city by city name and country code. Country code must be lowercase 2 letters: by, us, etc. 

```
curl GET "http://127.0.0.1:8080/weather?city=Gomel&country=by" -v
```

You can find a forecast for a city by city name and country code for any of next 7 days (0 - currect day, 6 - 7th day). 

```
curl GET "http://127.0.0.1:8080/weather?city=Gomel&country=by&day=2" -v
```


Configuration
-----
You can use environment variables: APIKEY(for open weather api), PORT, REDIS_HOST, REDIS_PORT:

```
REDIS_PORT=6789 go run main.go 
```
