# elastic-transcoder

To install node.js:
 ```bash
sudo-apt-get update
sudo apt-get install build-essential libssl-dev
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.0/install.sh | bash
nvm install v6.9.5
```
To execute the node.js server
 ```bash
cd DFS/app/
npm install
npm start
```

To compile both:
 ```bash
cd ~/goWorkspace/src/github.com/fercamp09/elastic-transcoder/
make

```
To execute the python client:
```bash
python rpc_clientpy.py <input-file.png> <output-file.jpg> <conversion-file-type> <priority> 
```

To execute the go client:
```bash
./client [create|cancel] <input-file.png> <output-file.jpg> <conversion-file-type> <priority>
```

To Execute the server:
```bash
./server 
```

In master branch, the files to process are previously located in all computers.

In prueba_final branch, is the python client that automatically uploads the file to the node.js server and the go worker downloads and process it.
