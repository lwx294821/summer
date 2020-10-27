#!/usr/bin/python3
# -*- coding: UTF-8 -*-
"""
@author:lishuwen
@file:Linux.py
@time:2020/08/21
"""

# coding: utf-8

import sys
import paramiko



# 定义一个类，表示一台远端linux主机
class Linux(object):
    # 通过IP, 用户名，密码，超时时间初始化一个远程Linux主机
    def __init__(self, ip, port, username, password, timeout=30):
        self.ip = ip
        self.port = port
        self.username = username
        self.password = password
        self.timeout = timeout
        # 连接失败的重试次数
        self.try_times = 3
        self.client=None

    # 调用该方法连接远程主机
    def connect(self):
        try:
            self.client = paramiko.SSHClient()
            self.client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            self.client.connect(hostname=self.ip, port=self.port, username=self.username, password=self.password)
            print('success connect the remote machine [host=%s]' % self.ip)
        except Exception as  e:
            print(
                'connect failed.in host[%s] user[%s] or pwd[%s] maybe wrong. ' % (
                    self.ip, self.username, self.password))
            sys.exit(-1)

    # 断开连接
    def close(self):
        if self.client:
           self.client.close()

    # 发送要执行的命令
    def send(self, cmd):
        try:
           stdin, stdout, stderr = self.client.exec_command(cmd,timeout=self.timeout,get_pty=True)
           result=[]
           status = stdout.channel.recv_exit_status()
           if status ==0:
               for i,l in  enumerate(stdout.readlines()):
                  result.append(str(l.replace("\r\n",'')))
           return status,result
        except Exception as e:
            return None,None
    
    def active(self,cmd):
        try:
           self.client.exec_command(cmd,timeout=self.timeout,get_pty=True)
           return 0
        except Exception as e:
            return None
        


if __name__ == '__main__':
    host = Linux('192.168.176.128', 22, 'root', '1',timeout=30)
    host.connect()
    result=host.active('docker stop elasticsearch')
    print(result)


