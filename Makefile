# This s how we want to name the binary output 
SERVER=server
CLIENT=client

# Builds the project 
build: 
	go build -o ${SERVER} rpc_server.go
	go build -o ${CLIENT} rpc_client.go

# Installs our project: copies binaries 
install: 
	go install 
  
# Cleans our project: deletes binaries 
clean: 
	if [ -f ${SERVER} ] ; then rm ${SERVER} ; fi
	if [ -f ${CLIENT} ] ; then rm ${CLIENT} ; fi

.PHONY: clean install 
