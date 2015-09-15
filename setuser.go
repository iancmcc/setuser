package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func abort(err error) {
	fmt.Fprintf(os.Stderr, "setuser: %s\n", err.Error())
	os.Exit(1)
}

func setupEnv(u *user.User) (err error) {
	var (
		uid int
		gid int
	)
	if uid, err = strconv.Atoi(u.Uid); err != nil {
		return
	}
	if gid, err = strconv.Atoi(u.Gid); err != nil {
		return
	}
	//if err = syscall.Setgroups([]int{gid}); err != nil {
	//	return
	//}
	_, _, e1 := syscall.RawSyscall(syscall.SYS_SETGID, uintptr(gid), 0, 0)
	if e1 != 0 {
		err = e1
		return
	}
	_, _, e1 = syscall.RawSyscall(syscall.SYS_SETUID, uintptr(uid), 0, 0)
	if e1 != 0 {
		err = e1
		return
	}
	os.Setenv("USER", u.Username)
	os.Setenv("HOME", u.HomeDir)
	os.Setenv("UID", u.Uid)
	return nil
}

func main() {
	var (
		u        *user.User
		err      error
		username string
		program  string
		command  []string
	)
	if len(os.Args[1:]) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s USERNAME COMMAND [args...]\n", os.Args[0])
		os.Exit(1)
	}
	username = os.Args[1]
	program = os.Args[2]
	command = append(command, os.Args[3:]...)
	if u, err = user.Lookup(username); err != nil {
		abort(err)
	}
	if program, err = exec.LookPath(program); err != nil {
		abort(err)
	}
	if err = setupEnv(u); err != nil {
		abort(err)
	}
	fmt.Println("Found binary at", program)
	fmt.Println("Args: ", command)
	if err = syscall.Exec(program, command, os.Environ()); err != nil {
		abort(err)
	}
}