#!/bin/sh

ab -n 100 -c 5 http://localhost:8000/hoge
ab -n 100 -c 5 http://localhost:8000/hello
ab -n 100 -c 5 http://localhost:8000/slow