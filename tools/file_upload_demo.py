# coding=utf-8
import base64

import requests


headers = {
    'Authorization': 'AEcuyNke3UYttRXrIr7BkQCditzCFGtMdZpFTedVDKJK',
    'Content-Type': 'application/json',
}
with open('/Users/Likeli/Desktop/1.png', 'rb') as f:
    ls_f = base64.b64encode(f.read())
json_data = {
    "msg": "",
    "type": 'image',
    "media_base64": ls_f.decode(),
    "extension_name": "png",
}
resp = requests.post('http://kf.api.7u1.cn/v1/app/dialog', headers=headers, json=json_data)
print(type(ls_f.decode()))
print(resp.status_code)
print(resp.json())
#
# imgdata = base64.b64decode(ls_f)
# with open('./2.png', 'wb') as f1:
#     f1.write(imgdata)
