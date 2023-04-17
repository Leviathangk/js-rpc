"""
    创建函数
    注意：已存在则覆写

    方法：
        create：指定域名的随机机器下创建函数
        create_by_uuid：指定 uuid 机器下创建函数
        create_more：指定域名的所有机器下创建函数
"""
import requests


def create(domain: str, func_name: str, func_body: str, timeout: int = None) -> dict:
    """
    在指定域名的机器上注入：如果有多个将会是随机

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
        url='http://127.0.0.1:8080/rpc/domain/create',
        data={
            'domain': domain,
            'funcName': func_name,
            'funcBody': func_body,
            'timeout': timeout,
        }
    )

    return resp.json()


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


def create_more(domain: str, func_name: str, func_body: str, timeout: int = None) -> dict:
    """
    在指定域名的所有机器上注入

    返回示例：内部消息同上
        {'msg': {'failed': [], 'failedTotal': 0, 'success': [{'msg': '覆写成功：hello', 'success': True, 'uuid': '7bf2e539-bb7b-4156-86b8-b66a478cd3a6'}], 'successTotal': 1, 'total': 1}, 'success': True}

    :param domain:
    :param func_name:
    :param func_body:
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/domain/more/create',
        data={
            'domain': domain,
            'funcName': func_name,
            'funcBody': func_body,
            'timeout': timeout,
        }
    )

    return resp.json()


if __name__ == '__main__':
    print(create(
        domain="newtab",
        func_name="hello",
        func_body="""
            console.log("执行完毕");
            sendResult("这是执行结果");
        """,
        timeout=5
    ))

    print(create_by_uuid(
        uuid="19f958a9-32c6-4a10-9ae7-196b888acde6",
        func_name="hello",
        func_body="console.log('你好')",
        timeout=5
    ))

    print(create_more(
        domain="newtab",
        func_name="hello",
        func_body="console.log('你好')",
        timeout=5
    ))
