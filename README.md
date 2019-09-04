# docker-go-starter

You can quick start for GO language development with docker.

It contain about Go v1.13 with Echo web flame-work and Gorm OR-Mapper with MySQL database.

You can edit go file which is automatically hot reload by fresh module.

Sample main.go is simple api just response json data. 

if you want to more variety for Dockerfile visit below.

https://github.com/docker-library/golang

## 1. Starting docker
Docke compose up and start API server. The sample endpoint is http://localhost:8080/helloWorld
When you edit the main.go, it will chang by fresh module has hot reload automatically.
```cassandraql
docker-compose up
Creating go ... done
Attaching to go
go     | 8:9:59 runner      | InitFolders
go     | 8:9:59 runner      | mkdir ./tmp
go     | 8:9:59 runner      | mkdir ./tmp: file exists
go     | 8:9:59 watcher     | Watching .
go     | 8:9:59 main        | Waiting (loop 1)...
go     | 8:9:59 main        | receiving first event /
go     | 8:9:59 main        | sleeping for 600 milliseconds
go     | 8:10:00 main        | flushing events
go     | 8:10:00 main        | Started! (5 Goroutines)
go     | 8:10:00 main        | remove tmp/runner-build-errors.log: no such file or directory
go     | 8:10:00 build       | Building...
go     | 8:10:04 runner      | Running...
go     | 8:10:04 main        | --------------------
go     | 8:10:04 main        | Waiting (loop 2)...
go     | 8:10:04 app         |
go     |    ____    __
go     |   / __/___/ /  ___
go     |  / _// __/ _ \/ _ \
go     | /___/\__/_//_/\___/ v3.3.10-dev
go     | High performance, minimalist Go web framework
go     | https://echo.labstack.com
go     | ____________________________________O/_______
go     |                                     O\

``` 
## 2. Browse by phpMyAdmin

http://localhost:8085/

You will see the test_database on your browser.

Check products table in the test_database when you post the request with Json parameter below.

## 3. hello world by Json post
It just looking your json parameter about your post request.
 
Call api from your local terminal.
```
curl   -X POST   http://localhost:8080/helloWorld   -H 'Content-Type: application/json'   -d '{"code": "helloWorld", "price": 999}'
```

## 4. Insert data by Json post
```cassandraql
curl -X POST http://localhost:8080/insert -H 'Content-Type: application/json' -d '{"PromotionCode": "foo","ContentUrl": "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", "StartedAt": "2019/09/03", "ExpiredAt": "2019/09/30"}'
```

## 5. Find data by Json Post
```cassandraql
curl   -X POST   http://localhost:8080/find   -H 'Content-Type: application/json'   -d '{"code": "foo"}'

{"code":200,"body":"555"}
```

## 6. Delete data by Json Post
```cassandraql
curl   -X POST   http://localhost:8080/delete   -H 'Content-Type: application/json'   -d '{"code": "foo"}'

{"code":200,"body":"foo"}
```

## 7. Login to go container
 
```cgo
docker exec -it go /bin/sh
```
