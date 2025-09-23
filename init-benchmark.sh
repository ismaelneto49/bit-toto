#!/bin/bash

for port in {5001..$1}; do
  # Other nodes run in background as servers
  go run src/main.go "$port" > logs/node_"$port".log 2>&1 &
done

# Node 5000 triggers searches
go run src/main.go 5000
