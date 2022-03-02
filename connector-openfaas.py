from kafka import KafkaConsumer
import json
from json import loads
import requests

url = 'http://10.110.204.22:8080/function/off'

connector = KafkaConsumer(
    'openfaas', bootstrap_servers=['129.114.25.12:9092'],
    auto_offset_reset='earliest', enable_auto_commit=True,
    group_id='0', value_deserializer=lambda x: loads(x.decode('utf-8')))

for message in connector:
    print("Inside Connector")
    message = message.value
    print(message)

    response = requests.post(url, data=message)

    print(response)