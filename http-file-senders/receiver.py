import requests
import json
import re

url = "http://localhost:3000/files"
idf = "589181bae9ea063bda4741aa" #el id del archivo, la metadata

http = requests.get(url +'/'+ idf)

headers = http.headers['content-disposition']
file_name = re.findall("filename=(.+)", headers)

with open('files/'+file_name[0], 'wb') as test: #direccion y nombre donde se guardara el archivo
    test.write(http.content)
