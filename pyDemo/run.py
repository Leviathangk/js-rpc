"""
    eval 函数

    方法：
        delete：指定域名的随机机器下删除函数
        delete_by_uuid：指定 uuid 机器下删除函数
        delete_more：指定域名的所有机器下删除函数
"""
import requests


def run(domain: str, func_name: str, timeout: int = None) -> dict:
    """
    在指定域名的机器上删除某个函数：如果有多个将会是随机

    返回示例：
        {"msg":"这是执行结果","success":true,"uuid":"7bf2e539-bb7b-4156-86b8-b66a478cd3a6"}

    :param domain:
    :param func_name:
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/domain/run',
        data={
            'domain': domain,
            'funcName': func_name,
            'timeout': timeout,
        }
    )

    return resp.json()


def run_by_uuid(uuid: str, func_name: str, timeout: int = None) -> dict:
    """
    在指定 uuid 的机器上删除某个函数

    返回示例：
        {"msg":"这是执行结果","success":true,"uuid":"7bf2e539-bb7b-4156-86b8-b66a478cd3a6"}

    :param domain:
    :param func_name:
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/uuid/run',
        data={
            'uuid': uuid,
            'funcName': func_name,
            'timeout': timeout,
        }
    )

    return resp.json()


def run_eval(domain: str, js: str, timeout: int = None) -> dict:
    """
    在指定域名的机器上删除某个函数：如果有多个将会是随机

    返回示例：
        {"msg":"","success":true,"uuid":"7bf2e539-bb7b-4156-86b8-b66a478cd3a6"}

    :param domain:
    :param js:
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/domain/eval',
        data={
            'domain': domain,
            'js': js,
            'timeout': timeout,
        }
    )

    return resp.json()


if __name__ == '__main__':
    print(run(
        domain="newtab",
        func_name="hello",
        timeout=5
    ))

    print(run_by_uuid(
        uuid="7bf2e539-bb7b-4156-86b8-b66a478cd3a6",
        func_name="hello",
        timeout=5
    ))

    print(run_eval(
        domain="newtab",
        js="console.log('1')",
        timeout=5
    ))
