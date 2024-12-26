package main

import (
	"encoding/json"
	"fake-ai-detective/config"
	"net"
)

var (
	openAICIDRList     []*net.IPNet
	cloudflareCIDRList []*net.IPNet
)

func init() {
	for _, cidr := range config.GetConfig().OpenAICIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		openAICIDRList = append(openAICIDRList, ipNet)
	}
	fetchAndUpdateCloudflareCIDR()
}

type ResponseInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CloudflareIPResponse struct {
	Result struct {
		Ipv4CIDRs []string `json:"ipv4_cidrs"`
		Ipv6CIDRs []string `json:"ipv6_cidrs"`
		Etag      string   `json:"etag"`
	} `json:"result"`
	Success  bool            `json:"success"`
	Errors   []*ResponseInfo `json:"errors"`
	Messages []*ResponseInfo `json:"messages"`
}

func fetchAndUpdateCloudflareCIDR() {
	url := "https://api.cloudflare.com/client/v4/ips"
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rsp := &CloudflareIPResponse{}
	err = json.NewDecoder(resp.Body).Decode(rsp)
	if err != nil {
		panic(err)
	}

	for _, cidr := range rsp.Result.Ipv4CIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		cloudflareCIDRList = append(cloudflareCIDRList, ipNet)
	}

	for _, cidr := range rsp.Result.Ipv6CIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		cloudflareCIDRList = append(cloudflareCIDRList, ipNet)
	}
}

func IsFromOpenAI(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return isIPInCIDR(ip, openAICIDRList)
}

func IsFromCloudflare(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return isIPInCIDR(ip, cloudflareCIDRList)
}

func isIPInCIDR(ip net.IP, cidrs []*net.IPNet) bool {
	for _, cidr := range cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}
