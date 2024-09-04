#!/bin/bash

docker build -t ascii . 
docker run -p :8082:8082 ascii
