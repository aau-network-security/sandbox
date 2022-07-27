// Copyright (c) 2018-2019 Aalborg University
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package dns

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aau-network-security/sandbox/store"
	"github.com/aau-network-security/sandbox/virtual"
	//"github.com/aau-network-security/defatt/sandbox"
	//"github.com/aau-network-security/sandbox2/store"
	//"github.com/aau-network-security/defatt/virtual"

	//"github.com/aau-network-security/defatt/virtual"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/aau-network-security/sandbox/virtual/docker"
	"github.com/rs/zerolog/log"
)

//var (
//	//go:embed Corefile.tmpl
//	Corefile embed.FS
//
//	//go:embed zonefile.tmpl
//	zonefile embed.FS
//)

//const (

//	coreFileContent = `. {
//    file zonefile
//    prometheus     # enable metrics
//    errors         # show errors
//    log            # enable query logs
//}
//`
//	zonePrefixContent = `$ORIGIN .
//@   3600 IN SOA sns.dns.icann.org. noc.dns.icann.org. (
//                2017042745 ; serial
//                7200       ; refresh (2 hours)
//                3600       ; retry (1 hour)
//                1209600    ; expire (2 weeks)
//                3600       ; minimum (1 hour)
//                )
//
//`
//)

type Server struct {
	cont     docker.Container
	corefile string
	zonefile string
	ipList   map[string]string
}

type Domains struct {
	Zonefile string
	URL      string
}

type RR struct {
	Type          string
	RData         string
	IPAddress     string
	Domain        string
	IPAddressMail string
	DomainMail    string
	IPAddressDC   string
	DomainDC      string
}

func createCorefile(domains Domains) string {
	var tpl bytes.Buffer

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir for CoreFile %v", err)
	}
	fullPathToTemplate := fmt.Sprintf("%s%s", dir, "/dnet/dns/Corefile.tmpl")

	tmpl := template.Must(template.ParseFiles(fullPathToTemplate))

	tmpl.Execute(&tpl, domains)
	return tpl.String()
}

func createZonefile(datas RR) string {

	var ztpl bytes.Buffer

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir for zonefile %v", err)
	}
	fullPathToTemplate := fmt.Sprintf("%s%s", dir, "/dnet/dns/zonefile.tmpl")

	tmpl := template.Must(template.ParseFiles(fullPathToTemplate))

	//tmpl := template.Must(template.ParseFiles("/home/ubuntu/vlad/sec03/defatt/dnet/dhcp/dhcpd.conf.tmpl"))

	tmpl.Execute(&ztpl, datas)
	return ztpl.String()
}

//func createMailZonefile(datas RR) string {
//
//	var mtpl bytes.Buffer
//
//	dir, err := os.Getwd() // get working directory
//	if err != nil {
//		log.Error().Msgf("Error getting the working dir for zonefile %v", err)
//	}
//	fullPathToTemplate := fmt.Sprintf("%s%s", dir, "/dnet/dns/zonemail.tmpl")
//
//	tmpl := template.Must(template.ParseFiles(fullPathToTemplate))
//
//	//tmpl := template.Must(template.ParseFiles("/home/ubuntu/vlad/sec03/defatt/dnet/dhcp/dhcpd.conf.tmpl"))
//
//	tmpl.Execute(&mtpl, datas)
//	return mtpl.String()
//}
//

func New(bridge string, ipList map[string]string, scenario store.Scenario, IPMail, IPdc string) (*Server, error) {

	var domains Domains
	var records RR

	domains.URL = scenario.FQDN
	//stripTLD := strings.SplitAfter(scenario.FQDN, ".")
	stripTLD := strings.Split(scenario.FQDN, ".")
	fmt.Printf("aici trebuie sa fie doar domaniul: %s\n", stripTLD[0])
	domains.Zonefile = stripTLD[0]

	c, err := ioutil.TempFile("", "Corefile")
	if err != nil {
		return nil, err
	}

	Corefile := c.Name()
	fmt.Printf("Asta este numele coreFile: %s\n", Corefile)
	CorefileStr := createCorefile(domains)

	_, err = c.WriteString(CorefileStr)
	if err != nil {
		return nil, err
	}

	for _, network := range ipList {

		ipAddrs := strings.TrimSuffix(network, ".0/24")
		ipAddrs = ipAddrs + ".2"
		records.IPAddress = ipAddrs

		for _, hosts := range scenario.Hosts {

			if hosts.Name == "mailserver" {
				records.IPAddressMail = IPMail
				//records.IPAddressMail = ConstructStaticIP(ipList,hosts.Networks,".3")
				records.DomainMail = hosts.DNS
			}
			if hosts.Name == "DCcon" {
				records.IPAddressDC = IPdc
				//records.IPAddressDC = sandbox.ConstructStaticIP(ipList,hosts.Networks,".251")
				records.DomainDC = hosts.DNS

			}

		}

	}

	records.Domain = scenario.FQDN

	z, err := ioutil.TempFile("", "zonefile")
	if err != nil {
		return nil, err
	}

	zonefile := z.Name()
	fmt.Printf("Asta este numele zonefile: %s\n", zonefile)

	zonefileStr := createZonefile(records)

	_, err = z.WriteString(zonefileStr)
	if err != nil {
		return nil, err
	}

	dir, err := os.Getwd() // get working directory
	if err != nil {
		log.Error().Msgf("Error getting the working dir for CoreFile %v", err)
	}

	//TODO: Make dynamic zonefile creation based on DNS

	cont := docker.NewContainer(docker.ContainerConfig{
		Image: "coredns/coredns:latest",
		Mounts: []string{
			fmt.Sprintf("%s:/Corefile", Corefile),
			fmt.Sprintf("%s:/root/db.%s", zonefile, domains.Zonefile),
			fmt.Sprintf("%s%s:/root/db.blue.monitor", dir, "/dnet/dns/db.blue.monitor"),
			//fmt.Sprintf("",dir ),
		},

		Resources: &docker.Resources{
			MemoryMB: 50,
			CPU:      0.3,
		},
		Cmd: []string{"--conf", "Corefile"},
		Labels: map[string]string{
			"nap-sandbox": bridge,
		},
	})

	return &Server{
		cont:     cont,
		corefile: Corefile,
		zonefile: zonefile,
	}, nil
}

func (s *Server) Container() docker.Container {
	return s.cont
}

func (s *Server) Run(ctx context.Context) error {
	return s.cont.Run(ctx)
}

//func (s *Server) Close(ctx context.Context) error {
//	if err := os.Remove(s.corefile); err != nil {
//		log.Warn().Msgf("error while removing DNS configuration file: %s", err)
//	}
//
//	if err := s.cont.Close(); err != nil {
//		log.Warn().Msgf("error while closing DNS container: %s", err)
//	}
//
//	return nil
//}

func (s *Server) Close() error {
	if err := os.Remove(s.corefile); err != nil {
		log.Warn().Msgf("error while removing DNS configuration file: %s", err)
	}

	if err := s.cont.Close(); err != nil {
		log.Warn().Msgf("error while closing DNS container: %s", err)
	}

	return nil
}

func (s *Server) Create(ctx context.Context) error {
	panic("implement me")
}

func (s *Server) Start(ctx context.Context) error {
	panic("implement me")
}

func (s *Server) Execute(ctx context.Context, i []string, s2 string) error {
	panic("implement me")
}

func (s *Server) Suspend(ctx context.Context) error {
	panic("implement me")
}

func (s *Server) Info() virtual.InstanceInfo {
	panic("implement me")
}

func (s *Server) Stop() error {
	return s.cont.Stop()
}
