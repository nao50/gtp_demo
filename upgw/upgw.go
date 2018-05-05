package main

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"./gtpv1"
	"golang.org/x/net/ipv4"
)

func main() {
	///////////////////////////////////////////////////////////////////////////////////////
	// common
	buffer := make([]byte, 1550)

	///////////////////////////////////////////////////////////////////////////////////////
	// S5:UPLINK:Recv:GTPv1Decap
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 2152,
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	///////////////////////////////////////////////////////////////////////////////////////
	// SGi:UPLINK:Send:RawSocket
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	defer syscall.Close(fd)

	///////////////////////////////////////////////////////////////////////////////////////
	// main loop
	fmt.Println("Starting raw server...")
	for {
		n, _, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			v1Packet := new(gtpv1.GTPV1)
			v1Packet.Parse(buffer[:n])

			ipheader, err := ipv4.ParseHeader(v1Packet.Data)
			if err != nil {
				fmt.Println("err: ", err)
			}

			addr := syscall.SockaddrInet4{
				Port: 0,
				Addr: [4]byte{ipheader.Dst.To4()[0], ipheader.Dst.To4()[1], ipheader.Dst.To4()[2], ipheader.Dst.To4()[3]},
			}

			err = syscall.Sendto(fd, v1Packet.Data, 0, &addr)
			if err != nil {
				log.Fatal("Sendto:", err)
			}
		}()

	}

}
