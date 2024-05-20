#!/bin/bash

export PGPASSWORD=postgres
psql -h 127.0.0.1 -p 5432 -U postgres -c 'drop database "consumer-api";'
psql -h 127.0.0.1 -p 5432 -U postgres -c 'create database "consumer-api";'
for i in `ls -1 *up.sql|sort`
    do echo ""
    echo $i
    psql -h 127.0.0.1 -p 5432 -U postgres -d "consumer-api" < $i
done
