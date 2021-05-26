# Sandbox

Version 0.0.1 of the sandbox, creates a virtual replica of interconnected networks and devices.

The starting point is the target VM, afterwards, networks are populated with additional virtual machine/Dockers.

### How to run?

To run this, you need to be on a Linux-based system. First, you will need to specify the folder where the platform can look
for *.ovafiles*. Update **ova-dir** from `config.yml` according to your needs.

> `config/config.yml`

Sandbox works on 4 parameters:

1. Tag -- Name of the experiment (e.g `test`); default comes with **test**; avoid **game**
2. vmsName -- .ovafile to used for population of networks; default value is **ubuntu.ova**
3. targetVM -- .ovafile used as target machine; default value is **ubuntu.ova**
4. networksNo -- number of networks to be created; default value is **3**

### Start sandbox

`go run main/main.go -tag=experiment -vmsName=windowsxp.ova -targetVM=wind7.ova -networksNo=5`

### Stop & Clean

`Ctrl+C`

`bash scripts/clean.sh`

### Requirements

- [**Go**](https://golang.org/doc/install)
- [**VirtualBox**](https://www.virtualbox.org/wiki/Linux_Downloads)
- [**Docker**](https://docs.docker.com/engine/install/ubuntu/)
- [**OpenVswitch**](https://www.openvswitch.org)
- [**ovs-docker**](https://github.com/aau-network-security/openvswitch/blob/master/scripts/ovs-docker) - this file
  should replace the one existing at
  `/usr/bin/ovs-docker`