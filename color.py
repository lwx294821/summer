#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:color.py
@time:2020/09/16
"""
import matplotlib.colors as c
import numpy as np
import hashlib

# 依据字符串标签文本取色值
def getColorByStr(str=''):
    b = str.encode('utf-8')
    css4 = c.CSS4_COLORS
    s = np.sum(list(map(int, b)))%len(css4)
    k=list(css4.keys())[s]
    return css4[k]

if __name__ == '__main__':
    print(getColorByStr('中国'))
