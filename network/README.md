
1. Configure the VM by using the virtualbox settings of that machine.
*Enable and set adapter 1 to NAT
*Enable and set adapter 2 to Bridged Adapter, allow VM.
*Save changes [OK]

2. Run and Enter the VM. Inside the VM:
*cd /etc/network/interfaces.d/
*If eth1.cfg file doesn't exist, create it by using: touch eth1.cfg
*vim /etc/network/interfaces.d/eth1.cfg # Modify file
* Use "i" key for insert, and write the following:
iface eth1 inet static
  address 192.168.12.13
  netmask 255.255.255.0
* Use ESC, and write ":wq" to save and quit.

3. To create the slaves, clone the VM.
4. In the new VM, "vim /etc/network/interfaces.d/eth1.cfg"
* Change the address "192.168.12.13" for 192.168.12.14 for slave 1

Do the same for slave 2
* Change the address "192.168.12.13" for "192.168.12.15" for slave 2

5. Copy ".erlang.cookie" file to "/var/lib/rabbitmq/" and "$HOME/"

6. Copy "etho.cfg" and "eht1.cfg" files to /etc/network/interfaces.d/

7. Copy "hosts" file to /etc/

8. Edit "/etc/hostname" file and write the name according to its ip. The hostname is the example file for master. 

192.168.12.13  master
192.168.12.14  slave1  
192.168.12.15  slave2
