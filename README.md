# nmapxmltocsv

Utility to convert NMAP XML output to a CSV file that looks like:

```
Hostname,IPAddress,Port,Protocol,Servicename,Servicestate
anshumanbhartiya.com,104.198.14.52,80,tcp,tcpwrapped,open
anshumanbhartiya.com,104.198.14.52,443,tcp,tcpwrapped,open
github.anshumanbhartiya.com,185.199.110.153,80,tcp,tcpwrapped,open
github.anshumanbhartiya.com,185.199.110.153,443,tcp,tcpwrapped,open
```

This is done using the `go-nmap` library of the [LAIR](https://github.com/lair-framework/go-nmap) framework.

## Running

`go run main.go -h`

`docker run -it abhartiya/tools_nmapxmltocsv -h`
