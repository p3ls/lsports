
# README

Prints information about used ports, processes, etc. It uses `fuser` and `netstat` under the hood. Only on Linux.

```
PORT    PROTO   PROCESS         USER    
22      tcp     sshd            root    
68      udp     dhcpcd          root    
546     udp     dhcpcd          root    
4569    udp     asterisk        asterisk
5038    tcp     asterisk        asterisk
5060    udp     asterisk        asterisk
5353    udp     avahi-daemon    avahi   
49853   udp     asterisk        asterisk
58467   udp     avahi-daemon    avahi   
59575   udp     asterisk        asterisk
59864   udp     avahi-daemon    avahi   

```

## Build

    go run main.go

