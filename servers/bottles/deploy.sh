#!/usr/bin/env bash
./build.sh

ssh -i ~/.ssh/finalPrivKey.pem ec2-user@18.222.243.235 << EOF
    
    docker rm -f bottles

    docker pull wecancodeit/bottles

    docker run -d \
    --name bottles \
    --network driftingNetwork \
    -e PORT=80 \
    wecancodeit/bottles
    
EOF