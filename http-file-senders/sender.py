import requests
import json
import re


url = "http://localhost:3000/files"
file_path = "sender.py" #la direccion del archivo a enviar
files = {'file': open(file_path, 'rb')}

r = requests.post(url, files=files)
json_data = json.loads(r.text)
file_id = json_data["_id"]
print (file_id + ' este es el id del archivo subido\n')
