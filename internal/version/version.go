package version

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
)

var (
	bi                  *debug.BuildInfo
	ErrCantGetBuildInfo = errors.New("can not get BuildInfo")
)

func GerProjectBuildInfo() (*debug.BuildInfo, error) {
	if bi == nil {
		var ok bool
		if bi, ok = debug.ReadBuildInfo(); !ok {
			return nil, ErrCantGetBuildInfo
		}
	}
	return bi, nil
}

func PrintVersionAndExit(string) error {
	bi, err := GerProjectBuildInfo()

	if err != nil {
		return err
	}

	fmt.Println(bi)

	os.Exit(0)

	return nil
}
