package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// EnsurePath is used to make sure a path exists
func EnsurePath(path string, dir bool) error {
	if !dir {
		path = filepath.Dir(path)
	}

	return os.MkdirAll(path, 0755)
}

func RandByte(n int) []byte {
	chars := []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u',
		'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
		'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_',
	}

	b := make([]byte, n)
	for i, l := 0, len(chars); i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		b[i] = chars[rand.Intn(l)]
	}

	return b
}

func GoId() (int, error) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GenName(name string) string {
	if name == "" {
		rs := string(RandByte(32))
		pid := int64(os.Getpid())
		wd := rs + "-" + strconv.FormatInt(pid, 10)

		goid, err := GoId()
		if err == nil {
			wd = wd + "-" + strconv.FormatInt(int64(goid), 10)
		}

		return GenMd5(wd)
	}

	return GenMd5(name)
}

func GenMd5(name string) string {
	hasher := md5.New()
	hasher.Write([]byte(name))

	return hex.EncodeToString(hasher.Sum(nil))
}

func GenHashCode(str string) uint {
	var hc uint = 0
	if str == "" {
		return hc
	}

	var n uint = uint(len(str))
	var i uint
	for i = 0; i < n; i++ {
		hc = hc ^ (uint(str[i]) << (uint(i) & 0xFF))
	}

	return hc
}

func GenHashDir(name string) []string {
	str := GenMd5(name)
	code := GenHashCode(str)
	hs := fmt.Sprintf("%06d", code)

	return []string{hs[0:2], hs[2:4], hs[4:6]}
}

func GetHashDir(dir, name string) (string, error) {
	hashDir := GenHashDir(name)

	d := path.Join(dir, strings.Join(hashDir, string(os.PathSeparator)))
	if err := EnsurePath(d, true); err != nil {
		if os.IsExist(err) {
			return d, nil
		}

		return "", err
	}

	return d, nil
}

func FindPrefixInStringArray(s string, arr []string) bool {
	for _, v := range arr {
		if strings.HasPrefix(s, v) {
			return true
		}
	}

	return false
}
