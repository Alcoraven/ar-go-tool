package artool

import "os"

func CreatePath(p string) error {
	if _, err := os.Stat(p); err != nil && os.IsNotExist(err) {
		e := os.MkdirAll(p, os.ModePerm)
		if e != nil {
			return e
		}
	}
	return nil
}
