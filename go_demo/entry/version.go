package main
import (
	"fmt"
	"os"
)

const (
	VERSION = "393fab9b435cd55e8527703214ffed30c1f8ea1b"
	BUILD_DATE = "2018-05-08 18:07"
	AUTHOR = "吴国福"
	BUILD_INFO = `Merge branch 'cherry-pick-109e140b' into 'master'`
)

func version() {
	fmt.Fprintf(os.Stderr, "version\t%s\n", VERSION)
	fmt.Fprintf(os.Stderr, "build\t%s\n", BUILD_DATE)
	fmt.Fprintf(os.Stderr, "author\t%s\n", AUTHOR)
	fmt.Fprintf(os.Stderr, "info\t%s\n", BUILD_INFO)
}

