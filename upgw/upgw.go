package main

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/naoyamaguchi/gtp/gtpv1"
	"golang.org/x/net/ipv4"
)

func main() {
	///////////////////////////////////////////////////////////////////////////////////////
	// common
	const proto = (syscall.ETH_P_IP<<8)&0xff00 | syscall.ETH_P_IP>>8
	uplinkBuffer := make([]byte, 1550)
	downlinkBuffer := make([]byte, 1550)

	///////////////////////////////////////////////////////////////////////////////////////
	// S5:DOWNLINK:Send:GTPv1Encap
	conn, err := net.Dial("udp4", "10.0.11.10:2152")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	///////////////////////////////////////////////////////////////////////////////////////
	// SGi:DOWNLINK:Recv:RawSocket
	recvSockFd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_DGRAM, proto)
	if err != nil {
		log.Fatal("recvSockFd: ", err)
	}
	defer syscall.Close(recvSockFd)

	recvIf, err := net.InterfaceByName("eth1")
	if err != nil {
		log.Fatal("interfacebyname: ", err)
	}

	recvSll := syscall.SockaddrLinklayer{
		Protocol: proto,
		Ifindex:  recvIf.Index,
	}

	if err := syscall.Bind(recvSockFd, &recvSll); err != nil {
		log.Fatal("bind: ", err)
	}

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
	fmt.Println("Starting pgw server...")
	for {
    go func() {
      ////////// up link //////////
      n, aa, err := udpConn.ReadFromUDP(uplinkBuffer)
      if err != nil {
        log.Fatalln(err)
      }
      // FOR DEBUG
      fmt.Printf("Recv udpaddr: %+v\n", aa)

      go func() {
        v1Packet := new(gtpv1.GTPV1)
        v1Packet.Parse(uplinkBuffer[:n])

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
    }()

    ////////// down link //////////
    go func() {
      n, addr, err := syscall.Recvfrom(recvSockFd, downlinkBuffer, 0)
      if err != nil {
        log.Fatalln(err)
      }
      // FOR DEBUG
      sa, _ := addr.(*syscall.SockaddrLinklayer)
      fmt.Printf("Recv SockaddrLinklayer: %+v\n", sa)

      go func() {
        g := &gtpv1.GTPV1{
          Version:                 1,
          ProtocolType:            1,
          Reserved:                0,
          ExtensionHeaderFlag:     0,
          SequenceNumberFlag:      1,
          N_PDUNumberFlag:         0,
          MessageType:             255,           // GPDU
          MessageLength:           uint16(n + 4), // よくわからん
          TEID:                    16879116,
          SequenceNumber:          65530,
          N_PDUNumber:             0,
          NextExtensionFeaderType: 0,
          Data: downlinkBuffer[:n],
        }
        msg := g.Marshal(downlinkBuffer[:n])

        _, err = conn.Write(msg)
        if err != nil {
          fmt.Println(err)
          return
        }
      }()
    }
	}

}
