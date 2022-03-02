from kafka import KafkaProducer
from json import dumps

import csv
import json


# Getting the data from Energy Data csv file
# Sending to kafka Broker so that FaaS function can listen from kafka brokers,
# in batches of 1000 records per batches

def csv_to_json():
    jsonArray = []
    csvFilePath = "energyData-Set2.csv"
    line_count = 1
    # read csv file
    rows = []
    with open(csvFilePath, encoding='utf-8') as csvf:
        # load csv file data using csv library's dictionary reader
        csvReader = csv.DictReader(csvf)

        # convert each csv row into python dict
        i = 0
        for row in csvReader:

            line_count = line_count + 1

            rows.append(row)
            if line_count % 1000 == 0:
                print('line count ', line_count)
                jsonArray = []
                # add this python dict to json array
                jsonArray.append(rows)
                rows = []

                producer = KafkaProducer(bootstrap_servers=['129.114.25.12:9092'],
                                         value_serializer=lambda x:
                                         dumps(x).encode('utf-8'))

                dict1 = {"docs": [{"msg": jsonArray[0]}]}
                jdata = json.dumps(dict1)
                ddata = json.loads(jdata)
                j2data = json.dumps(ddata)

                print(j2data)

                producer.send('openwhisk', value=j2data)
                producer.flush()

        print('line_count', line_count)


csv_to_json()
