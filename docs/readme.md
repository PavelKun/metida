Про JWT:
https://gist.github.com/zmts/802dc9c3510d79fd40f9dc38a12bccfc

Swagger:
go get -u github.com/swaggo/swag/cmd/swag
swag init
http://localhost:8080/swagger/index.html

Про профилирование:
http://localhost:8080/admin/debug/pprof/

go tool pprof metida http://localhost:8080/admin/debug/pprof/profile
```
>> top
>> web
>> list имя функции
```


одновременно с этим запустить апачи бенчмарк
ab -n 1000 -c 10 http://localhost:8080/show

Доступ:
sudo chmod -R 777 postgres-data/