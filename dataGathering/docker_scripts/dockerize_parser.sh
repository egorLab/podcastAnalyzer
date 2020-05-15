#! /bin/sh

docker build -t go-kubernetes .

docker tag go-kubernetes podcastalyzer/parser:1.0.0

docker login

docker push podcastalyzer/parser:1.0.0