package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// есть ли аргументы у вызова. без агрумента ошибка и выход.
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	PORT := arguments[1]
	// проверка доступности адреса на котором слушать порт
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	// слушать порт
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	// отложеное закрытие порта, когда он не будет болше использоваться
	defer connection.Close()
	// создание буфера для чтения порта длиной 1кб
	buffer := make([]byte, 1024)

        // настройка ипсет и иптаблес
        IpsetInit()

	// бесконечный цикл
	for {
		// читаем порт в буфер
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		// вывод в стдаут пришедшей строки. дебаг, можно удалить.
		fmt.Printf("%s\n", string(buffer[0:n]))

                // пример лога <190>Apr  6 04:52:38 nginx: 405 1.1.1.1
		// делим строку на 2 слова по слову нджинкс
		text := strings.Split(string(buffer[0:n]), "nginx: ")
		// делим вторую часть строки на слова по пробелу
		text = strings.Split(text[1], " ")
		status := text[0]
		ip := text[1]

		switch status { // если статус = то добавить в бан
		case "400":
			AddBannIp(ip)
		case "404":
			AddBannIp(ip)
		case "405":
			AddBannIp(ip)
		}
	}
}

func IpsetInit() {
        // создание списка autoban
        out, err := exec.Command("/sbin/ipset", "-!", "create",  "autoban", "hash:ip").CombinedOutput()
        if err != nil {
                fmt.Printf("%s", out)
                fmt.Println(err)
        }
        // очистка списка autoban, на всякий случай, если он был.
        out, err = exec.Command("/sbin/ipset", "-!", "flush",  "autoban").CombinedOutput()
        if err != nil {
                fmt.Printf("%s", out)
                fmt.Println(err)
        }
        // удаляю старое правило иптаблес. чтобы при рестарте не плодились правила
        // будет ошибка если его нет. поэтому ошибку игнорю.
        out, _ = exec.Command("/sbin/iptables", "-D", "INPUT", "-p", "tcp", "-m", "multiport", "--dports", "80,443", "-m", "set", "--match-set", "autoban", "src", "-j", "DROP").CombinedOutput()
        // добавляю новое.
        out, err = exec.Command("/sbin/iptables", "-I", "INPUT", "-p", "tcp", "-m", "multiport", "--dports", "80,443", "-m", "set", "--match-set", "autoban", "src", "-j", "DROP").CombinedOutput()
        if err != nil {
                fmt.Printf("%s", out)
                fmt.Println(err)
        }
}

// вызов ипсет добавление ип в список blocked-ips
func AddBannIp(BanIp string) {
	// сообщение о бане "адд ип". дебаг, можно удалить.
	fmt.Printf("add %s\n", BanIp)
	out, err := exec.Command("/sbin/ipset", "-!", "-A", "autoban", BanIp).CombinedOutput()
	if err != nil {
		fmt.Printf("%s", out)
		fmt.Println(err)
	}
}
