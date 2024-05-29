package main

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:


- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*




Так же требуется поддерживать функционал fork/exec-команд


Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).


*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Простейшее и универсальное решение: exec.Command("bash","-c",<то, что введено в командной строке>)

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type cmdType struct {
	cmd *exec.Cmd
	wg  sync.WaitGroup
	rc  io.ReadCloser
	pw  *io.PipeWriter
}

func execute(command string) {
	commands := strings.Split(command, "|")
	if len(commands) > 30 {
		fmt.Fprintln(os.Stderr, "too many pipes")
		os.Exit(1)
	}

	cmds := make([]cmdType, len(commands))

	for i, currentCommand := range commands {
		currentCommand = strings.TrimSpace(currentCommand)
		args := strings.Split(currentCommand, " ")
		cmds[i].cmd = exec.Command(args[0], args[1:]...)
		cmds[i].cmd.Stderr = os.Stderr
		if i == 0 {
			cmds[i].cmd.Stdin = os.Stdin
		}
		if i == len(commands)-1 {
			cmds[i].cmd.Stdout = os.Stdout
		}
		if i > 0 {
			switch cmds[i].cmd.Args[0] {
			case "pwd", "echo", "cd", "kill":
				cmds[i].cmd.Stdin, cmds[i-1].pw = io.Pipe() // Вывод предыдущей команды --> ввод текущей команды
				cmds[i-1].cmd.Stdout = cmds[i-1].pw         // Конструкция для возможности закрытия io.PipeWriter
			default:
				cmds[i].rc, _ = cmds[i-1].cmd.StdoutPipe() // Вывод предыдущей команды --> ввод текущей команды
				cmds[i].cmd.Stdin = cmds[i].rc             // Конструкция для возможности ручного закрытия io.ReadCloser
			}
		}
	}

	for i := range cmds {
		switch cmds[i].cmd.Args[0] {
		case "cd":
			cmds[i].cmd.Process = nil
			cmds[i].wg.Add(1)
			go func(x int) {
				defer cmds[x].wg.Done()
				defer func() {
					if cmds[x].rc != nil {
						cmds[x].rc.Close()
					}
					if cmds[x].pw != nil {
						cmds[x].pw.Close()
					}
				}()

				if len(cmds[x].cmd.Args) > 2 {
					fmt.Fprintln(os.Stderr, "Error in arguments:", cmds[x].cmd.Args)
					return
				}

				var home string

				if len(cmds[x].cmd.Args) > 1 {
					home = cmds[x].cmd.Args[1]
				} else {
					home = os.Getenv("HOME")
				}

				if err := os.Chdir(home); err != nil {
					fmt.Fprintln(os.Stderr, "Error chdir:", home, err.Error())
					return
				}
			}(i)
		case "pwd":
			cmds[i].cmd.Process = nil
			cmds[i].wg.Add(1)
			go func(x int) {
				defer cmds[x].wg.Done()
				defer func() {
					if cmds[x].rc != nil {
						cmds[x].rc.Close()
					}
					if cmds[x].pw != nil {
						cmds[x].pw.Close()
					}
				}()

				dir, _ := os.Getwd()
				fmt.Fprintln(cmds[x].cmd.Stdout, dir)

			}(i)
		case "echo":
			cmds[i].cmd.Process = nil
			cmds[i].wg.Add(1)
			go func(x int) {
				defer cmds[x].wg.Done()
				defer func() {
					if cmds[x].rc != nil {
						cmds[x].rc.Close()
					}
					if cmds[x].pw != nil {
						cmds[x].pw.Close()
					}
				}()

				if len(cmds[x].cmd.Args) > 1 {
					fmt.Fprintln(cmds[x].cmd.Stdout, strings.Join(cmds[x].cmd.Args[1:], " "))
					return
				}

				io.Copy(cmds[x].cmd.Stdout, cmds[x].cmd.Stdin)
			}(i)
		case "kill":
			cmds[i].cmd.Process = nil
			cmds[i].wg.Add(1)
			go func(x int) {
				defer cmds[x].wg.Done()
				defer func() {
					if cmds[x].rc != nil {
						cmds[x].rc.Close()
					}
					if cmds[x].pw != nil {
						cmds[x].pw.Close()
					}
				}()

				if len(cmds[x].cmd.Args) != 2 {
					fmt.Fprintln(os.Stderr, "There should be only one argument")
					return
				}

				pid, err := strconv.Atoi(cmds[x].cmd.Args[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, "Expected a number, you entered:", cmds[x].cmd.Args[1])
					return
				}

				proc, _ := os.FindProcess(pid) // проверяем наличие процесса на следующих строчках
				if err = proc.Signal(syscall.Signal(0)); err != nil {
					fmt.Fprintln(os.Stderr, "process doesn't exists PID:", pid)
					return
				}

				if err = proc.Kill(); err != nil {
					fmt.Fprintln(os.Stderr, "Error killing the process with PID:", pid, err.Error())
					return
				}
			}(i)
		default:
			if err := cmds[i].cmd.Start(); err != nil {
				fmt.Fprintln(os.Stderr, "Error with starting CMD:", cmds[i].cmd.Args, err.Error())
				return
			}
		}
	}

	for i := range cmds {
		if cmds[i].cmd.Process != nil {
			cmds[i].cmd.Wait()
		} else {
			cmds[i].wg.Wait()
		}
	}

}

func main() {
	const (
		prompt      = "myShell> "
		quitCommand = "\\quit"
	)

	var (
		cmdLine string
		err     error
	)

	for {
		fmt.Print(prompt)
		cmdLine, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Error input:", err.Error())
			os.Exit(1)
		}

		cmdLine = strings.TrimSpace(cmdLine)
		if len(cmdLine) == 0 {
			continue
		}
		if strings.HasPrefix(cmdLine, quitCommand) {
			fmt.Println("Good bye!")
			os.Exit(0)
		}
		execute(cmdLine)
	}
}
