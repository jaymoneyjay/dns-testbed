package experiment

import (
	"dns-testbed-go/testbed/component"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type SubqueryUnchained int

const (
	SubqueryBasic SubqueryUnchained = iota
	SubqueryCNAME
	SubqueryCNAME_QMIN
	SubqueryDNAME
	SubqueryDNAME_CVAL
	SubqueryDNAME_EXP
)

func (i SubqueryUnchained) String() string {
	return [...]string{"subquery", "subquery+CNAME", "subquery+CNAME+QMIN", "subquery+DNAME", "subquery+DNAME+CVAL", "subquery+DNAME+EXP"}[i]
}

type SlowDNS int

const (
	SlowDNS_Basic SlowDNS = iota
	SlowDNS_CNAME
	SlowDNS_CNAME_QMIN
)

func (i SlowDNS) String() string {
	return [...]string{"slowDNS", "slowDNS+CNAME", "slowDNS+CNAME+QMIN"}[i]
}

func saveQueryLog(target *component.Nameserver, fileName string, attackName string) error {
	logFile, err := os.Open(target.LogFile)
	if err != nil {
		return err
	}
	attackDirPath := fmt.Sprintf("results/logs/%s", attackName)
	if _, err := os.Stat(attackDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(attackDirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	destFile, err := os.Create(filepath.Join(attackDirPath, fileName))
	if err != nil {
		return err
	}
	_, err = io.Copy(destFile, logFile)
	if err != nil {
		return err
	}
	return nil
}
