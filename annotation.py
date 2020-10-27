#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:annotation.py
@time:2020/10/24
"""
from functools import wraps

from elasticsearch7 import Elasticsearch


def singleton(cls):
    _singleton = {}
    @wraps(cls)
    def wrapper(*args, **kwargs):
        if not _singleton.get(cls):
            _singleton[cls] = cls(*args, **kwargs)
        return _singleton[cls]
    return wrapper

@singleton
class ES(Elasticsearch):
    esHost = "192.168.176.128:9200"
    es = Elasticsearch(
        [esHost],
        sniff_on_start=False,
        sniff_on_connection_fail=True,
        sniffer_timeout=60,
        http_auth=('elastic', '123456')
    )
    def __init__(self, **kwargs):
        super().__init__(**kwargs)



if __name__ == '__main__':
    es1=ES()
    es2=ES()
    print(id(es1.es))
    print(id(es2.es))
