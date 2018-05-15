# GTP demo
GTP-U demonstration application in Golang
(GTP-C is comming soon ...)

## Getting Started
1. git clone & docker up.
```
$ git clone https://github.com/naoyamaguchi/gtp_demo.git
$ cd gtp_demo
$ docker-compose up -d
```
 NW diagram overview 

![nw-diagram](https://raw.githubusercontent.com/naoyamaguchi/gtp_demo/images/nwdiagram.png)

2. exec pseudo ue/gw container & configuration routing.
```
$ # ue container
$ docker-compose exec ue bash
$ sh init.sh

$ # gw container
$ docker-compose exec gw bash
$ sh init.sh
```
3. ping from ue to gw.
```
$ # gw container
$ tcpdump -i eth0

$ # ue container
$ ping 10.0.12.20
```
![nw-diagram-protocol](https://raw.githubusercontent.com/naoyamaguchi/gtp_demo/images/nwdiagram-protocol.png)

## Debug
If you are developing on a remote host, you can capture with local host Wireshark with the following command.
```
$ wireshark -k -i <(ssh -i ~/.ssh/<YOUR SSH PRIVATE KEY>.pem root@<REMOTE HOST IP ADDRESS> "tcpdump -U -n -w - -i any 'not port 22'") &
```


