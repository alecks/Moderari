#!/bin/sh
docker build -t moderari:latest .
docker run -d --rm --net host moderari:latest
