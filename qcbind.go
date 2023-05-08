/* qcbind - GPL - Copyright 2023 - condret */

package main

import "fmt"
import "os"
import "os/exec"
import "reflect"
import "unsafe"
import "golang.org/x/term"


func bytesToStr(b []byte) string {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := &reflect.StringHeader {
		Data: header.Data,
		Len: header.Len,
	}
	return *(*string)(unsafe.Pointer(strHeader))
}

func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func main () {
	if len (os.Args) < 3 {
		fmt.Printf (
			"Too few args :(\n" +
			"Usage: %s <path-to-image.qcow2> <path-to-nbd-device>\n", os.Args[0])
		os.Exit (-1)
	}
	var err error
	if _, err = os.Stat (os.Args[1]); err != nil {
		fmt.Printf ("Could no find %s :(\n", os.Args[1])
		os.Exit (-1)
	}
	if _, err = os.Stat (os.Args[2]); err != nil {
		fmt.Printf ("Could no find %s :(\n", os.Args[2])
		os.Exit (-1)
	}
	var image_arg string = fmt.Sprintf (
		"driver=qcow2,file.filename=%s,encrypt.format=luks,encrypt.key-secret=sec0", os.Args[1])
	var pw []byte
	pw, err = term.ReadPassword (0)
	if err != nil {
		fmt.Print (":(\n")
		os.Exit (-1)
	}
	// hack to clear the pw in keyarg after running qemu-nbd
	var keyarg_bytes []byte = append ([]byte("secret,id=sec0,data=")[:], pw[:]...)
	zeroBytes (pw)
	var keyarg string = bytesToStr (keyarg_bytes)
	var cmd *exec.Cmd = exec.Command ("qemu-nbd", "-c", os.Args[2], "--object", keyarg,
		"--image-opts", image_arg)
	cmd.Run ()
	zeroBytes (keyarg_bytes)
}
