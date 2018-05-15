# GTP demo
GTP-U demonstration application in Golang

## Getting Started
1. git clone & docker up
```
$ git clone https://github.com/naoyamaguchi/gtp_demo.git
$ cd gtp_demo
$ docker-compose up -d
$ 
```
2. exec pseudo ue/gw container & configuration routing
```
$ # ue container
$ docker-compose exec ue bash
$ sh init.sh

$ # gw container
$ docker-compose exec gw bash
$ sh init.sh
```
![nw-diagram](https://raw.githubusercontent.com/naoyamaguchi/gtp_demo/images/nwdiagram.png)
## Usage

## 