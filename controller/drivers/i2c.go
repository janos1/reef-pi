package drivers

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	delay    = 20
	slaveCmd = 0x0703 // Cmd to set slave address
	rdrwCmd  = 0x0707 // Cmd to read/write data together
	rd       = 0x0001
)

type message struct {
	addr  uint16
	flags uint16
	len   uint16
	buf   uintptr
}

type ioctlData struct {
	msgs uintptr
	nmsg uint32
}

type Bus struct {
	f  *os.File
	mu *sync.Mutex
}

func (b *Bus) send(cmd, addr uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.f.Fd(), cmd, addr); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func New() (*Bus, error) {
	f, err := os.OpenFile("/dev/i2c-1", os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	return &Bus{f: f, mu: new(sync.Mutex)}, nil
}
func (b *Bus) Close() error {
	return b.f.Close()
}

func (b *Bus) SetAddress(addr byte) error {
	return b.send(slaveCmd, uintptr(addr))
}

func (b *Bus) ReadBytes(addr byte, num int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.SetAddress(addr); err != nil {
		return []byte{0}, err
	}
	bytes := make([]byte, num)
	n, _ := b.f.Read(bytes)
	if n != num {
		return []byte{0}, fmt.Errorf("i2c: Unexpected number (%v) of bytes read", n)
	}
	return bytes, nil
}

func (b *Bus) WriteBytes(addr byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.SetAddress(addr); err != nil {
		return err
	}

	for i := range value {
		n, err := b.f.Write([]byte{value[i]})

		if n != 1 {
			return fmt.Errorf("i2c: Unexpected number (%v) of bytes written in WriteBytes", n)
		}
		if err != nil {
			return err
		}
		time.Sleep(delay * time.Millisecond)
	}
	return nil
}

func (b *Bus) ReadFromReg(addr, reg byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&value))

	var msgs [2]message
	msgs[0].addr = uint16(addr)
	msgs[0].flags = 0
	msgs[0].len = 1
	msgs[0].buf = uintptr(unsafe.Pointer(&reg))

	msgs[1].addr = uint16(addr)
	msgs[1].flags = rd
	msgs[1].len = uint16(len(value))
	msgs[1].buf = uintptr(unsafe.Pointer(hdrp.Data))

	var d ioctlData

	d.msgs = uintptr(unsafe.Pointer(&msgs))
	d.nmsg = 2
	if err := b.SetAddress(addr); err != nil {
		return err
	}
	return b.send(rdrwCmd, uintptr(unsafe.Pointer(&d)))
}

func (b *Bus) WriteToReg(addr, reg byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	outbuf := append([]byte{reg}, value...)

	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&outbuf))

	var msg message
	msg.addr = uint16(addr)
	msg.flags = 0
	msg.len = uint16(len(outbuf))
	msg.buf = uintptr(unsafe.Pointer(hdrp.Data))

	var d ioctlData
	d.msgs = uintptr(unsafe.Pointer(&msg))
	d.nmsg = 1

	if err := b.SetAddress(addr); err != nil {
		return err
	}
	return b.send(rdrwCmd, uintptr(unsafe.Pointer(&d)))
}
