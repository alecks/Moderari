#!/bin/sh
docker stop moderari
docker rmi moderari
docker build -t moderari .
docker run -d --rm --net host --name moderari moderari
