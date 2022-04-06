package main

import (
        "fmt"
        "net"
        "os"
        "os/exec"
        "strings"
)

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide a port number!")
                return
        }
        PORT := arguments[1]

        s, err := net.ResolveUDPAddr("udp4", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }

        connection, err := net.ListenUDP("udp4", s)
        if err != nil {
                fmt.Println(err)
                return
        }

        defer connection.Close()
        buffer := make([]byte, 1024)

        for {
                n, _, err := connection.ReadFromUDP(buffer)
                if err != nil {
                        fmt.Println(err)
                        return
                }
                fmt.Printf("%s\n", string(buffer[0:n]))

                text := strings.Split(string(buffer[0:n]), " ")
                if len(text) >= 7 {
                        if text[5] == "405" {
                                fmt.Printf("add %s\n", string(text[6]))
                                fmt.Println(AddBannIp(text[6]))
                        }
                }
        }
}

func AddBannIp(BanIp string) string {
        out, err := exec.Command("/sbin/ipset", "-A", "blocked-ips", BanIp).CombinedOutput()
        if err != nil {
                fmt.Println(err)
        }
        return string(out)
}
