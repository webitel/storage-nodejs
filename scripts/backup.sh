#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

t=$(date +%y%m%d%H%M%S)
[ -f $DIR/cc_schema.sql ] && cp $DIR/cc_schema.sql{,-$t}
docker exec -t mypg_db_1 /bin/sh -c 'pg_dump -U webitel --no-owner --clean --if-exists --schema=storage -s  webitel >/cc_schema.sql'
docker cp mypg_db_1:/cc_schema.sql $DIR/cc_schema.sql 
