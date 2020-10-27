#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:elasticsearch.py
@time:2020/08/14
"""

## 纠错

fuzzy = {
    "size": 0,
    "query": {
        "fuzzy": {
            "desc": {
                "value": "<KEYWORD>",
                "fuzziness": "AUTO",
                "prefix_length": 1,
                "max_expansions": 50,
                "transpositions": True,
                "boost": 1.0
            }
        }
    },
    "aggs": {
        "<AGG_NAME>": {
            "terms": {
                "field": "<FIELD>",
                "order": {
                    "max_score": "desc"
                }
            },
            "aggs": {
                "max_score": {
                    "max": {
                        "script": "_score"
                    }
                }
            }
        }
    }
}

## 建议
# mapping中的字段数据类型为completion
suggest = {
    "suggest": {
        "<FIELD_0>_suggest": {
            "prefix": "<KEYWORD>",
            "completion": {
                "field": "<COMPLETION_FIELD>",
                "size": 10
            }
        },
        "<FIELD_N>_suggest": {
            "prefix": "<KEYWORD>",
            "completion": {
                "field": "<COMPLETION_FIELD>",
                "size": 10,
                "fuzzy": {
                    "fuzziness": "AUTO",
                    "transposition": True,
                    "min_length": 3,
                    "prefix_length": 2,
                    "unicode_aware": True
                },
                "skip_duplicates": True
            }
        }
    }
}

## 聚合
composite = {
    "size": 0,
    "aggs": {
        "_buckets": {
            "composite": {
                "size": 20,
                "sources": [
                    {
                        "<FIELD>": {
                            "terms": {
                                "field": "<FIELD>",
                                "missing_bucket": True
                            }
                        }
                    }
                ]
            }
        }
    }
}

## 范围查询
range = {
    "query": {
        "bool": {
            "filter": {
                "range": {
                    "@timestamp": {
                        "time_zone": "+08:00",
                        "gte": "<START_TIME>",
                        "lte": "<END_TIME>",
                        "format": "yyyy-MM-dd HH:mm:ss.SSS"
                    }
                }
            },
            "must": [
                {
                    "term": {"<FIELD>": "<KEYWORD>"}
                }
            ]
        }
    }
}

# 判断消息体中是否包含某字段
exist_field = {
    "query": {
        "bool": {
            "filter": {
                "range": {
                    "@timestamp": {
                        "time_zone": "+08:00",
                        "gte": "<START_TIME>",
                        "lte": "<END_TIME>",
                        "format": "yyyy-MM-dd HH:mm:ss.SSS"
                    }
                }
            },
            "must": [
                {
                    "exists": {
                        "field": "<FIELD_NAME>"
                    }
                }
            ]
        }
    }
}

# 计算不重复数据数量
duplication = {
    "size": 0,
    "aggs": {
        "orgId": {
            "terms": {
                "field": "orgId.keyword"
            }
        },
        "count": {
            "cardinality": {
                "field": "orgId.keyword"
            }
        }
    }
}
groupby = {
    "aggs": {
        "id": {
            "terms": {
                "field": "id"
            },
            "aggs": {
                "having": {
                    "bucket_selector": {
                        "buckets_path": {
                            "view_count": "_count"
                        },
                        "script": "params.view_count < 1000"
                    }
                }
            }
        }
    }
}

# 分组后过滤每组数量


if __name__ == '__main__':
    pass
