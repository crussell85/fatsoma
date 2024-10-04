#!/bin/bash

docker-compose down -v
docker build -t ticket-api .
docker-compose up