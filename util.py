#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:util.py
@time:2020/08/12
"""
import time
import datetime
import pandas as pd

from common import color

NORMAL_FORMAT = '%Y-%m-%d %H:%M:%S.%f'
CST = '%a %b %d %H:%M:%S CST %Y'
TZ = '%Y-%m-%dT%H:%M:%SZ'

# 时间字符串转为时间对象
def getTimeByString(str):
    return datetime.datetime.strptime(str,NORMAL_FORMAT)


# CST格式时间转为本地时间
def cstToLocalTimeStr(str):
    to_format = '%Y-%m-%d %H:%M:%S'
    time_struct = time.strptime(str, CST)
    times = time.strftime(to_format, time_struct)
    return times

# TZ格式时间转为本地时间
def tzToLocalTimeStr(str):
    utc_date = datetime.datetime.strptime(str, TZ)
    local_date = utc_date + datetime.timedelta(hours=8)
    return datetime.datetime.strftime(local_date ,'%Y-%m-%d %H:%M:%S')

# 13位时间戳转为格式化时间
def timestamp_13_time(timestamp):
    return datetime.datetime.fromtimestamp(timestamp/1000.0).strftime(NORMAL_FORMAT)[:-3]

# 10位时间戳转为格式化时间
def timestamp_10_time(timestamp):
    return datetime.datetime.fromtimestamp(timestamp).strftime("%Y-%m-%d %H:%M:%S")

# 时间字符串转为时间戳
def str_to_timestamp(str_time, format_type='%Y-%m-%d %H:%M:%S'):
    time_array = time.strptime(str_time, format_type)
    return int(time.mktime(time_array)) * 1000

# 时间对象转为时间戳字符串
def getTimeStamp():
    datetime_object = datetime.datetime.now()
    now_timetuple = datetime_object.timetuple()
    now_second = time.mktime(now_timetuple)
    mow_millisecond = int(now_second * 1000 + datetime_object.microsecond / 1000)
    return mow_millisecond
# 时间对象转为格式化字符串
def getNowTimeStr(format='%Y-%m-%d %H:%M:%S'):
    return datetime.datetime.now().strftime(format)

# 时间字符串转为UNIX时间戳
def getUnixTimeStamp(str,format='%Y-%m-%d %H:%M:%S.%f'):
    time_strp = datetime.datetime.strptime(str,format)
    date_stamp = int(time.mktime(time_strp.timetuple()))
    date_microsecond = "%06d" % time_strp.microsecond
    timestamp ="%s.%s" % (date_stamp,date_microsecond[:3])
    return float(timestamp)




# 使用pandas批量将数据写入CSV文件
def to_csv(fileName=None, data=None, columns=None):
    columns = ["A", "B", "C"]
    pd.DataFrame(columns=columns).to_csv('kafka_msg.csv', sep=',', index=False, mode='w')
    items = []
    for i in range(5):
        items.append((str(i) + 'A', str(i) + 'B', str(i) + 'C'))
    df = pd.DataFrame.from_records(data=items)
    df.to_csv('kafka_msg.csv', sep=',', index=False, mode='a', header=False)



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
    # c = cstToLocalTimeStr('Tue Nov 06 00:00:00 CST 2018')
    # print(c)
    #
    # t = tzToLocalTimeStr('2019-07-26T08:20:54Z')
    # print(t)

    # print(timestamp_13_time(1597630795798))
    print(getUnixTimeStamp("2018-11-22 17:20:23.010"))

