#!/bin/bash




#request sudo....
#if [[ $UID != 0 ]]; then
#    echo "Please run this script with sudo:"
#    echo "sudo $0 $*"
#    exit 1
#fi

while getopts "b:" arg; do
  case $arg in
    b) bridge=$OPTARG;;
  esac
done

sudo kill $(ps -e | pgrep tcpdump);

sudo ovs-vsctl del-br $bridge
sudo ip link del ${bridge}_vlan10
sudo ip link del ${bridge}_vlan20
sudo ip link del ${bridge}_vlan30
sudo ip link del ${bridge}_AllBlue
sudo ip link del ${bridge}_monitoring
sudo ip link del ${bridge}_DMZ
sudo ip link del ${bridge}_special

VBoxManage list runningvms | awk '/'$bridge'/ {print $1}' | xargs -I vmid VBoxManage controlvm vmid poweroff
VBoxManage list vms | awk '/'$bridge'/ {print $1}' | xargs -I vmid VBoxManage unregistervm --delete vmid
VBoxManage list runningvms | awk '/'$bridge'/ {print $2}' | xargs -I vmid VBoxManage controlvm vmid poweroff
VBoxManage list vms | awk '/'$bridge'/ {print $2}' | xargs -I vmid VBoxManage unregistervm --delete vmid

VBoxManage list vms | awk '/'target-*'/ {print $1}' | xargs -I vmid VBoxManage unregistervm --delete vmid


rm -rf ~/VirtualBox\ VMs/nap-$bridge-*
rm -rf ~/VirtualBox\ VMs/$bridge-*

#while read -r line; do
 #   vm=$(echo $line | cut -d ' ' -f 2)
  #  echo $vm
  #  vboxmanage controlvm $vm poweroff
   # vboxmanage unregistervm $vm --delete
#done <<< "$VMS"



# Remove all docker containers that have a UUID as name
#docker ps -a --format '{{.Names}}' | grep -E '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' | xargs docker rm -f

docker kill $(docker ps -q -a -f "label=sandbox-${bridge}")

docker rm $(docker ps -q -a -f "label=sandbox-${bridge}")
#docker system prune --filter "label=FTPstatic"

# Prune entire docker
docker system prune -f --filter "label=sandbox-${bridge}"
#docker system prune --filter "label=FTPstatic"

# Prune volumes
docker volume prune -f --filter "label=sandbox-${bridge}"
