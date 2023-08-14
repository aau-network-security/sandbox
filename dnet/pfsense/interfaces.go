package pfsense

import (
	"github.com/rs/zerolog/log"
	"strconv"
)

//}

type AvailableInterfaces struct {
	Mac      string `json:"mac"`
	Up       bool   `json:"up"`
	Ipaddr   string `json:"ipaddr"`
	Friendly string `json:"friendly"`
	Dmesg    string `json:"dmesg"`
	InUse    string `json:"in_use"`
}

type Interfaces struct {
	Status  string                         `json:"status"`
	Code    int                            `json:"code"`
	Return  int                            `json:"return"`
	Message string                         `json:"message"`
	Data    map[string]AvailableInterfaces `json:"data"`
}

func (c *Client) GetInterfaceList() (*Interfaces, error) {

	payload := map[string]string{}
	req, err := c.prepareRequest("GET", "/interface/available", payload)

	if err != nil {
		return nil, err
	}

	resp := Interfaces{}
	if err := c.sendRequest(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ConfigureInterface struct {
	AdvDhcpConfigAdvanced         bool     `json:"adv_dhcp_config_advanced"`
	AdvDhcpConfigFileOverride     bool     `json:"adv_dhcp_config_file_override"`
	AdvDhcpConfigFileOverrideFile string   `json:"adv_dhcp_config_file_override_file"`
	AdvDhcpOptionModifiers        string   `json:"adv_dhcp_option_modifiers"`
	AdvDhcpPtBackoffCutoff        int      `json:"adv_dhcp_pt_backoff_cutoff"`
	AdvDhcpPtInitialInterval      int      `json:"adv_dhcp_pt_initial_interval"`
	AdvDhcpPtReboot               int      `json:"adv_dhcp_pt_reboot"`
	AdvDhcpPtRetry                int      `json:"adv_dhcp_pt_retry"`
	AdvDhcpPtSelectTimeout        int      `json:"adv_dhcp_pt_select_timeout"`
	AdvDhcpPtTimeout              int      `json:"adv_dhcp_pt_timeout"`
	AdvDhcpRequestOptions         string   `json:"adv_dhcp_request_options"`
	AdvDhcpRequiredOptions        string   `json:"adv_dhcp_required_options"`
	AdvDhcpSendOptions            string   `json:"adv_dhcp_send_options"`
	AliasAddress                  string   `json:"alias-address"`
	AliasSubnet                   int      `json:"alias-subnet"`
	Apply                         bool     `json:"apply"`
	Blockbogons                   bool     `json:"blockbogons"`
	Blockpriv                     bool     `json:"blockpriv"`
	Descr                         string   `json:"descr"`
	Dhcpcvpt                      int      `json:"dhcpcvpt"`
	Dhcphostname                  string   `json:"dhcphostname"`
	Dhcprejectfrom                []string `json:"dhcprejectfrom"`
	Dhcpvlanenable                bool     `json:"dhcpvlanenable"`
	Enable                        bool     `json:"enable"`
	Gateway                       string   `json:"gateway"`
	Gateway6Rd                    string   `json:"gateway-6rd"`
	Gatewayv6                     string   `json:"gatewayv6"`
	Id                            string   `json:"id"`
	If                            string   `json:"if"`
	Ipaddr                        string   `json:"ipaddr"`
	Ipaddrv6                      string   `json:"ipaddrv6"`
	Ipv6Usev4Iface                bool     `json:"ipv6usev4iface"`
	Media                         string   `json:"media"`
	Mss                           string   `json:"mss"`
	Mtu                           int      `json:"mtu"`
	Prefix6Rd                     string   `json:"prefix-6rd"`
	Prefix6RdV4Plen               int      `json:"prefix-6rd-v4plen"`
	Spoofmac                      string   `json:"spoofmac"`
	Subnet                        int      `json:"subnet"`
	Subnetv6                      string   `json:"subnetv6"`
	Track6Interface               string   `json:"track6-interface"`
	Track6PrefixIdHex             int      `json:"track6-prefix-id-hex"`
	Type                          string   `json:"type"`
	Type6                         string   `json:"type6"`
}

func (c *Client) CreateInterface(v ConfigureInterface) (*ConfigureInterface, error) {

	payload := map[string]string{
		"apply":  strconv.FormatBool(v.Apply),
		"type":   v.Type,
		"descr":  v.Descr,
		"if":     v.If,
		"id":     v.Id,
		"enable": strconv.FormatBool(v.Enable),
	}
	req, err := c.prepareRequest("POST", "/interface", payload)

	if err != nil {
		return nil, err
	}

	res := v
	if err := c.sendRequest(req, &res); err != nil {
		log.Err(err).Msgf("Error is: %v, now is %s", err, v.Id)
		return nil, err
	}

	return &res, nil
}

//curl -X POST \
//'http://192.168.1.253/api/v1/interface/apply' \
//-H 'accept: application/json' \
//-H 'Content-Type: application/json' \
//-d '{
//"async": false
//}'

func (c *Client) ApplyInterface(v ConfigureInterface) (*ConfigureInterface, error) {

	payload := map[string]string{
		"async": strconv.FormatBool(false),
	}
	req, err := c.prepareRequest("POST", "/interface", payload)

	if err != nil {
		return nil, err
	}

	res := ConfigureInterface{}
	if err := c.sendRequest(req, &res); err != nil {
		log.Err(err).Msgf("Error is: %v, now is %s", err, v.Id)
		return nil, err
	}

	return &res, nil

}
