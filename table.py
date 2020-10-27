#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:table.py
@time:2020/10/24
"""


from common.annotation import ES
from common.util import setLabelColor


def table(data):
    es = ES()
    body= {
            "size":50,
            "query": {
                 "bool": {
                     "must": [
                        {
                            "match": {
                                'label': data['label']
                            }
                        }
                    ],
                     "term":{
                         "type":{
                             "value":data['type']
                         }
                     }
                }
            },
            "sort": [
                {"_id": "desc"}
            ]
        }
    if len(data['after'])>0:
        body["search_after"]=data['after']
    res=es.search(index="workspace",body=body)
    after = ''
    hits=[]
    for h in res['hits']['hits']:
        h['_source']['id'] = h['_id']
        if isinstance(h['_source'], dict):
            if dict(h['_source']).get('label') is not None:
                labels = setLabelColor(h['_source']['label'])
                h['_source']['label'] = labels
        hits.append((h['_source']))
        after = h['sort']
    return after, hits

if __name__ == '__main__':
    print("Hello")