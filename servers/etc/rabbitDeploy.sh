#!/usr/bin/env bash

ssh -i ~/.ssh/MyPrivKey.pem ec2-user@18.217.182.145 << EOF
    
    docker rm -f rabbitmq

    docker run -d --name rabbitmq \
    --network driftingNetwork \
    rabbitmq:3-management

    exit

EOF