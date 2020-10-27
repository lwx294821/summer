#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:doc_controller.py
@time:2020/09/14
"""
import os
import time
import traceback
import numpy as np
import pandas as pd
from elasticsearch7 import Elasticsearch

from common import tracer, util, color

logger = tracer.log
esHost = "192.168.176.128:9200"
es = Elasticsearch(
    [esHost],
    sniff_on_start=False,
    sniff_on_connection_fail=True,
    sniffer_timeout=60,
    http_auth=('elastic', '123456')
)

'''
记录文档到Elasticsearch
'''


# 创建索引
def create_index():
    indices = ['material']
    mappings = {
        "settings": {
            "number_of_replicas": 0,
            "number_of_shards": 1
        },
        'mappings': {
            'dynamic': True
        }
    }
    for index in indices:
        isIn = es.indices.exists(index=index)
        if isIn:
            es.indices.delete(index=index)
        time.sleep(5)
        es.indices.create(index=index, body=mappings)


def remove_doc(node):
    es.delete(index=node['index'], id=node['id'])


def index_doc():
    data = {
        'zone': '河西',
        'department': '智能电网服务中心',
        'project': '用采2.0',
        'type': '虚机登录信息',
        'name': '服务机器分布',
        'index': 'material',
        'item': '备忘录',
        'status': '1',
        'timestamp': util.getNowTimeStr(),
        'label': 'machine service Test',

        'log': '172.17.37.41 | root | Ft041!test | java -jar /home/NettyTest_jar/NettyTest.jar',
        'k8s': '26.47.131.167 | root | 123qwe!@# | /root/txwg/devops/activeReport/deploy.yaml'

    }

    # res = es.index(index='material', body=data)

    args = {
        'name': '/opt/kafka_2.12-2.4.1/bin/kafka-topics.sh --zookeeper {BROKER_SERVER}:2181 --delete --topic {TOPIC_NAME}',
        'value': "",
        'default': "",
        'version': "2.4.1",
        'label': "Kafka Topic Delete"
    }
    res = es.index(index='open', body=args)


# 滚动查询所有
def get_scroll_doc(index='open', label='', search_after=None):
    if search_after is None or len(search_after)==0:
        body = {
            "query": {
                # "bool": {
                #     "must": [
                #         {
                #             "match": {
                #                 'label': label
                #             }
                #         }
                #     ]
                # }
                "match_all":{}
            },
            "sort": [
                {"_id": "desc"}
            ]
        }
        res = es.search(index=index, body=body)
        print(res)
    else:
        body = {
            "query": {
                # "bool": {
                #     "must": [
                #         {
                #             "match": {
                #                 'label': label
                #             }
                #         }
                #     ]
                # }
                "match_all":{}
            },
            "search_after": [search_after],
            "sort": [
                {"_id": "desc"}
            ]
        }
        res = es.search(index='open', body=body)
        print(res)
    hits=[]
    after=''
    for h in res['hits']['hits']:
        h['_source']['id'] = h['_id']
        if isinstance(h['_source'], dict):
            if dict(h['_source']).get('label') is not None:
                labels = setLabelColor(h['_source']['label'])
                h['_source']['label'] = labels
        hits.append((h['_source']))
        after = h['sort'][0]
    return after, hits


def search_doc(node):
    label = node['label']
    if len(label) == 0 or len(node['type']) ==0:
        body = {
            "query": {
                "match_all": {}
            }
        }
    else:
        body = {
            "query": {
                "bool": {
                    "must": [
                        {
                            "match": {
                                'label': label
                            }
                        },
                        {
                            'term':{
                                'type':{
                                    'value':node['type']
                                }
                            }
                        }
                    ]
                }
            }
        }

    res = es.search(index=node['index'], body=body)
    content = []
    for h in res['hits']['hits']:
        if isinstance(h['_source'], dict):
            if dict(h['_source']).get('label') is not None:
                labels = setLabelColor(h['_source']['label'])
                h['_source']['label'] = labels
                h['_source']['id'] = h['_id']
                content.append(h['_source'])
    return content


def clear_all(index=""):
    body = {
        "query": {
            "match_all": {
            }
        }
    }
    es.delete_by_query(index=index, body=body)


def get_all_doc(index):
    scroll_id = None
    while True:
        try:
            scroll_id, hits = scroll_doc(scroll_id, index)
            if len(hits) <= 0:
                break
            else:
                yield hits
        except Exception as e:
            print(e)


def scroll_doc(scroll_id, index):
    body = {
        "query": {
            "match_all": {}
        }
    }
    if scroll_id is None:
        res = es.search(index=index, body=body, scroll='5m')
    else:
        scroll_body = {
            "scroll": "5m",
            "scroll_id": scroll_id
        }
        res = es.scroll(body=scroll_body)
    hits = []
    for h in res['hits']['hits']:
        hits.append(h['_source'])
    return res['_scroll_id'], hits


'''
备份所有索引中的所有数据
'''
def backup_all(indexes=None):
    os.remove('es_all.csv')
    if isinstance(indexes, list):
        for index in indexes:
            for h in get_all_doc(index):
                df = pd.DataFrame.from_records(data=h)
                df.to_csv('es_all.csv', index=False, mode='a', header=False)


def setLabelColor(label=''):
    labels = label.split()
    label_obj = []
    for l in labels:
        label_obj.append(
            {
                'name': l,
                'color': color.getColorByStr(l)
            }
        )
    return label_obj


if __name__ == '__main__':
    # es.delete(index='open',id='qZVC3HQBBVoOzBQK6huF')
    #backup_all(indexes=['open', 'material', 'compute','project'])
    data = {
        'name': 'Elastisearch索引生命周期策略',
        'address': '江苏省电力公司客户服务大厦奥体大街9号',
        'type': 'problem',
        'label': 'es ilm',
        'level': 'L1',
        'status': 'plan',
        'timestamp': util.getNowTimeStr(),
        'required':'hot:20%,replicas 1,cold:80%,replicas 0'

    }

    body = {
        "query": {
            # "bool": {
            #     "must": [
            #         {
            #             "match": {
            #                 'label': label
            #             }
            #         }
            #     ]
            # }
            "match_all": {}
        },
        "sort": [
            {"_id": "desc"}
        ]
    }
    index_doc()
