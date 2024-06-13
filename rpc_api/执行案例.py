"""
    创建函数
    注意：已存在则覆写

    方法：
        create：指定域名的随机机器下创建函数
        create_by_uuid：指定 uuid 机器下创建函数
        create_more：指定域名的所有机器下创建函数
"""
import json

import requests


def create_by_uuid(uuid: str, func_name: str, func_body: str, timeout: int = None) -> dict:
    """
    在指定 uuid 的机器上注入

    返回示例：
        {"msg":"注入成功","success":true,"uuid":"d1de0b19-6f7a-4afd-b290-bc415a2ab307"}
        {"msg":"覆写成功：hello","success":true,"uuid":"7bf2e539-bb7b-4156-86b8-b66a478cd3a6"}
        {"msg":"缺失字段","success":false}

    :param domain:
    :param func_name:
    :param func_body:
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/uuid/create',
        data={
            'uuid': uuid,
            'funcName': func_name,
            'funcBody': func_body,
            'timeout': timeout,
        }
    )

    return resp.json()


def run_by_uuid(uuid: str, func_name: str, args: any = None, timeout: int = None) -> dict:
    """
    在指定 uuid 的机器上删除某个函数

    返回示例：
        {"msg":"这是执行结果","success":true,"uuid":"7bf2e539-bb7b-4156-86b8-b66a478cd3a6"}

    :param uuid:
    :param func_name:
    :param args: 需要传的参数
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/uuid/run',
        data={
            'uuid': uuid,
            'funcName': func_name,
            'timeout': timeout,
            'args': json.dumps(args)
        }
    )

    return resp.json()


if __name__ == '__main__':
    print(create_by_uuid(
        uuid="6c8f6bd0-d9ff-4ec4-9e06-1da48f37d56f",
        func_name="getData",
        func_body="""
            sendResult(window.getData(recvJson.token))
        """,
        timeout=5
    ))

    print(run_by_uuid(
        uuid="6c8f6bd0-d9ff-4ec4-9e06-1da48f37d56f",
        func_name="getData",
        args={
            "token": "225e3eed9e9a4b5aac9bc496b2952625"
        },
        timeout=5,
    ))
