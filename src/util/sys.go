package util

import (
	"os"
	"strconv"
	"unsafe"
)

const intSize = int(unsafe.Sizeof(0))

var bs *[intSize]byte

func init() {
	i := 0x1
	bs = (*[intSize]byte)(unsafe.Pointer(&i))
}

func IsBigEndian() bool {
	return !IsLittleEndian()
}

func IsLittleEndian() bool {
	return bs[0] == 0
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func HostName() (hostname string) {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "UNKNOWN"
	}
	return
}

func GetEnvInt(name string, def int) int {
	env, ok := os.LookupEnv(name)
	if ok {
		i64, err := strconv.ParseInt(env, 10, 0)
		if err != nil {
			return def
		}
		return int(i64)
	}
	return def
}

func GetEnvString(name string, def string) string {
	env, ok := os.LookupEnv(name)
	if ok {
		return env
	}
	return def
}
