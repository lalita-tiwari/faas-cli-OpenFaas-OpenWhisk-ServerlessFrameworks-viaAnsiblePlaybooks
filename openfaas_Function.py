import requests
import json

def handle(args):

    print(args)
    print('updated $$$$$$$$$$$$$$$$$$$$$')
    url = "http://129.114.25.124:32302/openfaas/_bulk_docs"

    headers = {'Content-type': 'application/json', 'Accept': 'text/plain'}

    response = requests.post(url, data=args, auth=('admin', 'admin'), headers=headers)

    dict = {"pass":"success"}
    return dict