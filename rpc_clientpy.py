#!/usr/bin/env python
import pika
import uuid
import tasks.task_pb2 as pb
import requests
import json
import re
#1from tasks.task_pb2 import Response
import sys


class imageRpcClient(object):
    def __init__(self):
        self.connection = pika.BlockingConnection(pika.ConnectionParameters(
                host='master'))

        self.channel = self.connection.channel()

        result = self.channel.queue_declare(exclusive=True)
        self.callback_queue = result.method.queue

        self.channel.basic_consume(self.on_response, no_ack=True,
                                   queue=self.callback_queue)



    def on_response(self, ch, method, props, body):
        if self.corr_id == props.correlation_id:
            self.response = body

    def call(self, n):
        self.response = None
        self.corr_id = str(uuid.uuid4())
        self.channel.basic_publish(exchange='',
                                   routing_key='rpc_queue2',
                                   properties=pika.BasicProperties(
                                         reply_to = self.callback_queue,
                                         correlation_id = self.corr_id,
                                         priority=n.priority,
                                         ),
                                         body=n.SerializeToString()
                                   )
        while self.response is None:
            self.connection.process_data_events()
        return self.response


url = "http://localhost:3000/files"
#file_path = "sender.py" #la direccion del archivo a enviar

task= pb.Task()
if len(sys.argv) < 4:
  print ("Please supply an input and output filename e.g. go run rpc_client.go input.jpg output.jpg jpg 1")
  sys.exit(-1)
elif len(sys.argv) == 4:
    #files = {'file': open(sys.argv[1], 'rb')}
    #r = requests.post(url, files=files)
    #json_data = json.loads(r.text)
    #file_id = json_data["_id"]

    task.priority=0
    task.filename = sys.argv[1] #file_path
    task.new_name = sys.argv[2]
    #task.FileId = sys.argv[3]
else:
    #files = {'file': open(sys.argv[1], 'rb')}
    #r = requests.post(url, files=files)
    #json_data = json.loads(r.text)
    #file_id = json_data["_id"]
    
    task.filename = sys.argv[1]
    task.new_name = sys.argv[2]
    task.priority = int(sys.argv[4])
    #task.FileId = file_id

image_rpc = imageRpcClient()
print(" [x] Requesting image")
response = image_rpc.call(task)
res=pb.Response()
res.ParseFromString(response)
file_loc=res.file_location
print(" [.] Image processed found in %s" % file_loc)
