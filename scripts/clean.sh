#!/bin/bash


# todo: will be modified with more dynamic way of cleaning stuff

#request sudo....
#if [[ $UID != 0 ]]; then
#    echo "Please run this script with sudo:"
#    echo "sudo $0 $*"
#    exit 1
#fi

sudo ovs-vsctl del-br test
sudo ip tuntap del tap0 mode tap
sudo ip tuntap del tap10 mode tap
sudo ip tuntap del tap20 mode tap
sudo ip tuntap del tap30 mode tap
sudo ip tuntap del tap40 mode tap
sudo ip tuntap del tap50 mode tap
sudo ip tuntap del tap60 mode tap
sudo ip tuntap del vlan110 mode tap
sudo ip tuntap del vlan220 mode tap
sudo ip tuntap del vlan330 mode tap
sudo ip tuntap del vlan440 mode tap
sudo ip tuntap del vlan550 mode tap
sudo ip tuntap del mon10 mode tap
sudo ip tuntap del ALLblue mode tap


VBoxManage list runningvms | awk '/sandbox/ {print $1}' | xargs -I vmid VBoxManage controlvm vmid poweroff
VBoxManage list vms | awk '/sandbox/ {print $2}' | xargs -I vmid VBoxManage unregistervm --delete vmid


rm -rf ~/VirtualBox\ VMs/sandbox*

#while read -r line; do
 #   vm=$(echo $line | cut -d ' ' -f 2)
  #  echo $vm
  #  vboxmanage controlvm $vm poweroff
   # vboxmanage unregistervm $vm --delete
#done <<< "$VMS"



# Remove all docker containers that have a UUID as name
#docker ps -a --format '{{.Names}}' | grep -E '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' | xargs docker rm -f

docker kill $(docker ps -q -a -f "label=sandbox")

docker rm $(docker ps -q -a -f "label=sandbox")




# Remove all macvlan networks
docker network rm $(docker network ls -q -f "label=sandbox")

# Prune entire docker
docker system prune --filter "label=sandbox"

# Prune volumes
docker volume prune --filter "label=sandbox"