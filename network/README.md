# Configure VMs cluster in the same Network
## Configure bridged network in Virtualbox

1. Configure the VM by editing the virtualbox settings of that machine
  1.Enable and set adapter 1 to NAT
  1.Enable and set adapter 2 to Bridged Adapter, allow VM
  1.Save changes [OK]

2. Run and Enter the VM. Inside the VM:
   Alternatively, you can copy the "etho.cfg" and "eht1.cfg" files from this directory into `/etc/network/interfaces.d/` or follow the following steps:

  1. `cd /etc/network/interfaces.d/`. 
  2. If "eth1.cfg" file doesn't exist, create it by using: `touch eth1.cfg`.  
  3. `sudo vim /etc/network/interfaces.d/eth1.cfg # Modify file`
  4. Inside VIM, Use "i" key for insert, and write the following:
   iface eth1 inet static
   address 192.168.12.13
   netmask 255.255.255.0
  5. Use ESC, and write ":wq" to save and quit

3. To create the slaves, clone the VM.
4. In the new VM, "sudo vim /etc/network/interfaces.d/eth1.cfg"
* Change the address "192.168.12.13" for 192.168.12.14 for slave 1

1. Repeat steps 3 and 4, for slave 2
* Change the address "192.168.12.13" for "192.168.12.15" for slave 2

## Configure rabbitmq to be able to create cluster.
1. Copy ".erlang.cookie" file to "/var/lib/rabbitmq/" and "$HOME/"
2. Copy "hosts" file to /etc/

## Configure rabbitmq to be able to communicate with guest user accross VM.
1. Edit "/etc/hostname" file and write the name according to its ip. The hostname file in this directory is the example file for master. 

### IPS and hostnames:
- 192.168.12.13  master
- 192.168.12.14  slave1  
- 192.168.12.15  slave2

## Create cluster.
Inside a slave node, input:
If the rabbitmq-server is running, stop it using `sudo rabbitmqctl stop`. Start a rabbitmq detached node server by inputting: `sudo rabbitmq-server -detached`. The master node also has to be detached, so also do this to the master.

Then form the cluster by stopping the node, joining master cluster and starting the node:
```
sudo rabbitmqctl stop_app
sudo rabbitmqctl join_cluster vagrant@master
rabbitmqctl start_app
```
