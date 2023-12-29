# MaP1058-socket-client
Go socket client capturing bio-signals sent from [MaP1058](https://wp.santeku-map.com/%e3%82%bd%e3%83%95%e3%83%88%e3%82%a6%e3%82%a7%e3%82%a2/%e5%8f%8e%e9%8c%b2%e3%83%bb%e8%a7%a3%e6%9e%90%e3%82%bd%e3%83%95%e3%83%88%e3%82%a6%e3%82%a7%e3%82%a2/map1058/)

## Usage
Short example of usage:
```golang
package your_package

import (
    "time"
    "github.com/ISDL-dev/MaP1058-socket-client/client"
)

func YourFunc() {
    conf := client.Config{
        ServerIP: "192.168.10.101",
        SaveDir:  "./output",
    }
    c, _ := client.NewClient(conf)
    c.Start(time.Minute * 5)
}
```
The above code will create a client receiving bio-signals from MaP1058 and save them in the `./output` directory for 5 minutes.  
Full example is at [cmd/example/main.go](cmd/example/main.go).