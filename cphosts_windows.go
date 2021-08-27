package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	psikuta "psikuta"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

var hostsPath = os.Getenv("windir") + "\\system32\\drivers\\etc\\hosts"

func main() {
	fmt.Println("Program aktualizuje hostsy na podstawie serwera HP. Wersja: 2021.08 build: 27")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var username, pass, destHost, jakaPodsiec string
	username, pass, destHost, jakaPodsiec = psikuta.HPcred()
	sshConfig, err := auth.PasswordKey(username, pass, ssh.InsecureIgnoreHostKey())
	checkError(err)
	cpFromHP(sshConfig, destHost, jakaPodsiec)
}

func cpToHP(sshConfig ssh.ClientConfig, destHost string) {
	scpClient := scp.NewClient(destHost, &sshConfig)

	err := scpClient.Connect()
	checkError(err)

	localfileData, err := os.Open("go.mod")
	checkError(err)

	scpClient.CopyFile(localfileData, "go.mod", "0655")

	defer scpClient.Session.Close()
	defer localfileData.Close()
}

func cpFromHP(sshConfig ssh.ClientConfig, destHost string, jakaPodsiec string) {
	scpServer := scp.NewClient(destHost, &sshConfig)
	remoteDir := "/etc/hosts"

	err := scpServer.Connect()
	checkError(err)

	localfileData, err := os.OpenFile("hosts.lst", os.O_RDWR|os.O_CREATE, 0666)
	checkError(err)

	scpServer.CopyFromRemote(localfileData, remoteDir)

	defer scpServer.Session.Close()
	defer localfileData.Close()

	kopiaZapasowa()
	filtrujNaszeIP("hosts.lst", jakaPodsiec)
	removeDuplicates(hostsPath)
}

func filtrujNaszeIP(file string, jakaPodsiec string) {
	search := flag.String(jakaPodsiec, jakaPodsiec, jakaPodsiec)
	data, _ := ioutil.ReadFile(file)

	for _, line := range strings.Split(string(data), "\n") {
		if strings.Index(line, *search) > -1 {
			fmt.Printf("W pliku %s znalazłem: %s\n", file, line)
			f, err := os.OpenFile(hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
			checkError(err)
			defer f.Close()
			f.WriteString(line + "\n")
			checkError(err)
		}
	}
}

func removeDuplicates(file string) {
	line, _ := ioutil.ReadFile(file)
	strLine := string(line)
	lines := strings.Split(strLine, "\n")

	found := make(map[string]bool)
	j := 0
	for i, x := range lines {
		if !found[x] {
			found[x] = true
			(lines)[j] = (lines)[i]
			j++
		}
	}
	lines = (lines)[:j]
	f, err := os.OpenFile(hostsPath, os.O_WRONLY, 0666)
	checkError(err)
	err = f.Truncate(0)
	_, err = f.Seek(0, 0)
	checkError(err)
	f.Close()
	f, err = os.OpenFile(hostsPath, os.O_WRONLY, 0666)
	checkError(err)

	fmt.Println(lines)
	for e := range lines {
		f.Write([]byte(lines[e] + "\n"))
	}
	defer f.Close()
	os.Remove("hosts.lst")
	os.Remove("hosts.tmp")
}

func kopiaZapasowa() {
	bytesRead, err := ioutil.ReadFile(hostsPath)
	checkError(err)
	const layout = "2006-01-02"
	t := time.Now()
	err = ioutil.WriteFile(hostsPath+"."+t.Format(layout), bytesRead, 0644)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
