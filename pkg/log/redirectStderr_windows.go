// Copyright (c) 2022 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

//go:build windows
// +build windows

package log

import (
	"log"
	"os"
	"syscall"
)

var (
	_kernel32         = syscall.MustLoadDLL("_kernel32.dll")
	_procSetStdHandle = _kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(_procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) error {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
		return err
	}
	// SetStdHandle does not affect prior references to stderr
	os.Stderr = f
	return nil
}
