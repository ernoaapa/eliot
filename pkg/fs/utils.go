package fs

import "os"

// FileExist returns true if file exist in given path
func FileExist(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && !stat.IsDir()
}

// DirExist returns true if directory exist in given path
func DirExist(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && stat.IsDir()
}
