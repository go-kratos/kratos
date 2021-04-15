# Run the examples

## gin
1、Compile and execute the server code:
```shell
$ go run gin/main.go
```
2、From a different terminal, access the api to see the output:
```shell
$ curl http://127.0.0.1:8000/home
Hello Gin!
```

## mux
1、Compile and execute the server code:
```shell
$ go run mux/main.go
```
2、From a different terminal, access the api to see the output:
```shell
$ curl http://127.0.0.1:8000/home
Hello Gorilla Mux!
```

## static
1、Compile and execute the server code:
```shell
$ go run static/main.go
```
2、Access the following url with your browser:
```
http://127.0.0.1:8000/assets
```