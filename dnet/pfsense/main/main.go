package main

import (
	"github.com/aau-network-security/sandbox/dnet/pfsense"
)

//import "pfsense"
//curl -H "Authorization: 61646d696e 8d32c42e4d9225f3051ad57153cd42b9" -X GET  'http://192.168.0.203/api/v1/user'

//curl -H "Authorization: 61646d696e 8d32c42e4d9225f3051ad57153cd42b9" -X POST \
//'http://192.168.1.253/api/v1/interface/apply' \
//-H 'accept: application/json' \
//-H 'Content-Type: application/json' \
//-d '{
//"async": false
//}'

func main() {

	clint := pfsense.NewClient("http://192.168.1.253", "61646d696e", "8d32c42e4d9225f3051ad57153cd42b9")
	//NewClient("http://192.168.0.203", "61646d696e", "8d32c42e4d9225f3051ad57153cd42b9")

	//fmt.Println("doamne nu stiu nimic")
	//useL, err := clint.GetUserList()
	//if err != nil {
	//	log.Error().Msgf("Error getting the working dir %v", err)
	//}
	//fmt.Println(useL.Data)
	//
	//for _, datum := range useL.Data {
	//	fmt.Printf("%s; %s;%s; %s",datum.Username,datum.Password,datum.Descr,datum.AuthorizedKeys)
	//}

	//intf, err := clint.GetInterfaceList()
	//
	//if err != nil {
	//	log.Error().Msgf("Error getting the interfaces %v", err)
	//
	//}
	//fmt.Println(intf)
	//
	//for _,value := range intf.Data {
	//	fmt.Println(value)
	//}
	//fmt.Println(s)

	//TODO:plm, merge dar pune descr gresit si nu ia enable. trebuie verificat cum formeaza Jsonul
	// Pwp ma duc la concert
	v := pfsense.ConfigureInterface{
		AliasAddress: "",
		AliasSubnet:  0,
		Apply:        true,
		Descr:        "pula",
		Enable:       true,
		If:           "em2",
		Id:           "opt1",
		Type:         "dhcp",
		Type6:        "dhcp6",
	}
	//v := pfsense.ConfigureInterface{
	//
	//	Apply:             true,
	//	Enable:            true,
	//		If:                "em2",
	//		Id: "opt1",
	//	Type:              "dhcp",
	//}
	clint.CreateInterface(v)

}
