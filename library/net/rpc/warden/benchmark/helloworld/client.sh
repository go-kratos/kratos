#!/bin/bash
go build -o client greeter_client.go
echo size 100 concurrent 30
./client -s 100 -c 30
echo size 1000 concurrent 30
./client -s 1000 -c 30
echo size 10000 concurrent 30
./client -s 10000 -c 30
echo size 100 concurrent 300
./client -s 100 -c 300
echo size 1000 concurrent 300
./client -s 1000 -c 300
echo size 10000 concurrent 300
./client -s 10000 -c 300
rm client