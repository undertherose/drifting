#!/usr/bin/env bash

#builds docker container
docker build -t wecancodeit/bottles .

docker push wecancodeit/bottles

