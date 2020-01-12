#!/usr/bin/env bash

wrk -t1 -c1 -d10s http://127.0.0.1:8080
wrk -t1 -c10 -d10s http://127.0.0.1:8080
wrk -t1 -c100 -d10s http://127.0.0.1:8080
wrk -t1 -c200 -d10s http://127.0.0.1:8080
wrk -t1 -c300 -d10s http://127.0.0.1:8080