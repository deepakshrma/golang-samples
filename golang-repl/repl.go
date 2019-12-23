package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/micro/go-micro/util/log"
)

func printMore() {
	fmt.Print("\n(0/exit)-> ")
}

func executeAsync(cmd []string) {
	ex := exec.Command(cmd[0], cmd[1:]...)
	cmdReader, _ := ex.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf(scanner.Text())
		}
	}()
	ex.Start()
	err := ex.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
func execute(cmd []string) {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	output := string(out[:])
	fmt.Println(output)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("## Go REPL ##")
	fmt.Println("---------------------")
	exit := false
	spaceReg := regexp.MustCompile(`\s`)

	for !exit {
		printMore()
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		cmd = strings.TrimSpace(cmd)
		if strings.Compare("exit", cmd) == 0 || strings.Compare("0", cmd) == 0 {
			fmt.Println("Thanks, Have a good day!")
			exit = true
		} else {
			executeAsync(spaceReg.Split(cmd, -1))
		}
	}
}
