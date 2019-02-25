#!/usr/bin/env bash
./build.sh

docker push koolkids441/finalgateway

export TLSCERT=/etc/letsencrypt/live/api.iqueue.zubinchopra.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.iqueue.zubinchopra.me/privkey.pem
export REDISADDR=redisServer:6379
export MONGOADDR=mongo:27017
export MYSQL_ROOT_PASSWORD="sqldbpassword"
export DSN="root:%s@tcp\(finalsqldb:3306\)/auth"
export COURSESADDR=courses:80
export MQNAME=rabbitmq
export MQADDR=rabbitmq:5672
export FAQADDR=faq:80 


ssh -i ~/.ssh/finalPrivKey.pem ec2-user@18.222.243.235 << EOF

    docker rm -f redisServer
    docker rm -f driftingServer
    
    docker pull wecancodeit/gateway

    docker run -d --name redisServer \
    --network driftingNetwork \
    redis

    docker run -d --name driftingServer --network driftingNetwork -p 443:443 -e REDISADDR=$REDISADDR -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e DSN=$DSN -e COURSESADDR=$COURSESADDR -e FAQADDR=$FAQADDR -e MONGOADDR=$MONGOADDR -e MQADDR=$MQADDR -e MQNAME=$MQNAME wecancodeit/gateway
    
    exit
EOF

