package main

import "github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library"

func main() {
	proxy := exhibition_proxy_library.Proxy{
		Port: 8080,
		SettingsPath: "proxy-settings.json",
	}
	proxy.StartServer()

}
