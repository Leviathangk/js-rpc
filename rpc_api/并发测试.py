import sys
from pathlib import Path

sys.path.append(str(Path(__file__).absolute().parent.parent))

from rpc_api import create, run
from concurrent.futures import ThreadPoolExecutor

# print(create.create(
#     domain="newtab",
#     func_name="getData",
#     func_body="""
#         sendResult(new Date().getTime())
#     """,
#     timeout=5
# ))

def fire():
    print(run.run(
        domain="newtab",
        func_name="getData"
    ))

threads = 20
with ThreadPoolExecutor(max_workers=threads) as t:
    for _ in range(threads):
        t.submit(fire)