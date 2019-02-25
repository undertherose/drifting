#!/usr/bin/env bash

ssh -i ~/.ssh/finalPrivKey.pem ec2-user@18.222.243.235 << EOF
    
    docker rm -f rabbitmq

    docker run -d --name rabbitmq \
    --network driftingNetwork \
    rabbitmq:3-management

    exit

EOF