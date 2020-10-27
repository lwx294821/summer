#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:storage.py
@time:2020/10/24
"""
from elasticsearch7 import Elasticsearch

esHost = "192.168.176.128:9200"
es = Elasticsearch(
    [esHost],
    sniff_on_start=False,
    sniff_on_connection_fail=True,
    sniffer_timeout=60,
    http_auth=('elastic', '123456')
)

def add(index,data):
    try:
        es.index(index=index, body=data)
    except Exception:
        return ""
    return "Congratulation."

if __name__ == '__main__':
    pass