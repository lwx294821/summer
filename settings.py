#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:settings.py
@time:2020/10/24
"""
from common.annotation import ES


def esIndex(data):
    es = ES()
    index = data['index']
    if data['action'] == "0":
        mappings = {
            "settings": {
                "number_of_replicas": 0,
                "number_of_shards": 1
            },
            'mappings': {
                'dynamic': True
            }
        }
        es.indices.create(index=index, body=mappings)
    if data['action'] == "-1":
        isIn = es.indices.exists(index=index)
        if isIn:
            es.indices.delete(index=index)


if __name__ == '__main__':
    pass
