#!/bin/sh
docker rmi moderari
docker build -t moderari .
docker run -d --rm --net host moderari
