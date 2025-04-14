package main

import "github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy"

func main() {
	proxy := exhibition_proxy.Proxy{
		Port: 8080,
		SettingsPath: "proxy-settings.json",
	}
	proxy.StartServer()

}
