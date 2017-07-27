#!/usr/bin/env python
# -*- coding: utf-8 -*-
import sys
import os
from redis import Redis
from psql import PSQL

# ENV = "local"
# ENV = "dev"
ENV = "dev"
if os.environ.get('ENV') != "":
    ENV = os.environ.get('ENV')

REDIS_URL = {
    "local": "redis://localhost:6379/0",
    "dev": "redis://:myredis@11.239.191.199:6379/0",
    "prod": "redis://:myredis@localhost:6379/0"
}

PG_URL = {
    "local": "postgres://localhost:5432/app",
    "dev": "postgres://app:appapp@11.239.191.199:5432/app",
    "prod": "postgres://localhost:5432/app"
}

redis = Redis.from_url(REDIS_URL[ENV])
pg = PSQL(PG_URL[ENV])

wdj_done_apk_sql = """SELECT apk FROM android WHERE wdj = 2;"""
wdj_todo_chan = "wdj:app:todo"
wdj_done_chan = "wdj:app:done"
wdj_fail_chan = "wdj:app:fail"


def run_sql(chan, sql):
    todo = pg.fetch_column(sql)
    n = len(todo)
    bufSize = 10000
    buf = []
    for i, apk in enumerate(todo):
        buf.append(apk)
        if i % bufSize == 0:
            redis.lpush(chan, *buf)
            print("Commit: %d / %d" % (i, n))
            buf = []
    if len(buf) > 0:
        redis.lpush(chan, *buf)


# example:
# python driver.py "wdj:app:todo" "SELECT apk FROM android WHERE wdj = 1 LIMIT 100;"
# python driver.py "wdj:search:todo" "SELECT name FROM wdj_app LIMIT 100;"


if __name__ == '__main__':

    if len(sys.argv) == 3:
        chan = sys.argv[1]
        sql = sys.argv[2]
        print("[Env:%s], [Chan:%s] [CMD:%s]" % (ENV, chan, sql))
        run_sql(chan, sql)
