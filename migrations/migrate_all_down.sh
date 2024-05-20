#!/bin/bash

export PGPASSWORD=postgres
for i in `ls -1 *down.sql|sort`
    do echo ""
    echo $i
    psql -h 127.0.0.1 -p 5432 -U postgres -d "consumer-api" < $i
done
