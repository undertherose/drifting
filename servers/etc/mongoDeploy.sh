#!/usr/bin/env bash

ssh -i ~/.ssh/finalPrivKey.pem ec2-user@18.222.243.235 << EOF

    docker rm -f finalmongo

    docker run -d --name finalmongo \
    --network queueNetwork \
    mongo

    exit
EOF
