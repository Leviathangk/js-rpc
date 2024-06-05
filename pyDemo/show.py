from typing import List

import requests


def show_domain_clients(domain: str, timeout: int = None) -> List[str]:
    """
    返回指定域名下的所有客户端的 uuid

    返回示例：
        {"msg":{"clients":["15254d28-f001-49b2-941f-ac34ff7a2386"],"total":1},"success":true}

    :param domain:
    :param func_name:
    :param func_body:
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/show_domain_clients',
        data={
            'domain': domain,
            'timeout': timeout,
        }
    )
    print(resp.text)
    print(resp)

    content = resp.json()

    if content['success'] and content['msg']['clients']:
        return content['msg']['clients']

    return []


def show_client_functions(uuid: str, timeout: int = None) -> List[str]:
    """
    返回指定客户端的所有函数名

    返回示例：
        {"msg":{"functions":[],"total":0},"success":true,"uuid":"15254d28-f001-49b2-941f-ac34ff7a2386"}

    :param domain:
    :param func_name:
    :param func_body:
    :param timeout: timeout 为 None 即无限等待
    :return:
    """
    resp = requests.post(
        url='http://127.0.0.1:8080/rpc/show_client_functions',
        data={
            'uuid': uuid,
            'timeout': timeout,
        }
    )
    print(resp.text)
    print(resp)

    content = resp.json()

    if content['success'] and content['msg']['total']!=0:
        return content['msg']['functions']

    return []


if __name__ == '__main__':
    show_domain_clients(
        domain="newtab",
        timeout=5
    )

    show_client_functions(
        uuid="0d44260b-88e4-4340-8d04-cd462b42ea35",
        timeout=5
    )
