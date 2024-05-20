#!/bin/bash
pwd
ls
export PGPASSWORD=postgres
for i in `ls -1 *up.sql|sort`
    do echo ""
    echo $i
    psql -h postgres -p 5432 -U postgres -d "consumer-api" < $i
done