
# OpenWhisk function

# rename the file as __main__.py


import requests


import requests
import json

def main(args):


    url = "http://129.114.25.124:32302/openwhisk/_bulk_docs"

    headers = {'Content-type': 'application/json', 'Accept': 'text/plain'}

    response = requests.post(url, data=json.dumps(args), auth=('admin', 'admin'), headers=headers)
    print(response)
    dict = {"pass":"success"}
    return dict
