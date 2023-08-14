package controller

import "fmt"

//TCPdump *TCPdump

type TCPdump struct {
	c *NetController
}

//sudo tcpdump -i any -nn -w webserver.pcap port 80

func (dump *TCPdump) DumpTraffic(intf, file string) error {
	cmds := []string{"-i", fmt.Sprintf("%s", intf), "-nn", "-w", fmt.Sprintf("%s.pcap", file)}
	//_, err := ipc.exec(fmt.Sprintf("tuntap del %s mode %s", tap, mode))
	_, err := dump.exec(cmds...)

	return err
}

//func (dump *IPService) StopRecording(intf, file string) error {
//	cmds := []string{"tcpdump", "-i", intf, "-nn", "-w", fmt.Sprintf("%s.pcap", file)}
//	//_, err := ipc.exec(fmt.Sprintf("tuntap del %s mode %s", tap, mode))
//	_, err := dump.exec(cmds...)
//
//	return err
//}

func (dump *TCPdump) exec(args ...string) ([]byte, error) {
	return dump.c.exec("tcpdump", args...)
}
