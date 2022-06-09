package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Abs(Path string) (string, error) {
	if !filepath.IsAbs(Path) {
		p, err := filepath.Abs(Path)
		if err != nil {
			return "", err
		}
		Path = p
	}
	return Path, nil
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func GetFilesInDir(dir string) (res []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		res = append(res, f.Name())
	}
	return
}

func CreateDir(dir string) (string, error) {
	if exist := FileExists(dir); exist {
		return dir, nil
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	return dir, nil
}

func CreateFile(dir string, filename string, source io.Reader) (string, error) {
	dirPath, err := CreateDir(dir)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dirPath, filename)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, source)
	if err != nil {
		return "", err
	}
	return path, nil
}
