//go:build !js

package service

import (
	"gioui.org/app"
	"github.com/mearaj/bhagad-house-booking/common/alog"
	"os"
	"path/filepath"
)

func (s *service) openDatabase() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			if err != nil {
				alog.Logger().Errorln(err)
			}
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		dirPath, err := app.DataDir()
		if err != nil {
			return
		}
		dirPath = filepath.Join(dirPath, DBPathCfgDir)
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0700)
			if err != nil {
				return
			}
		}
		dbFullName := filepath.Join(dirPath, DBPathFileName)
		if _, err = os.Stat(dbFullName); os.IsNotExist(err) {
			var file *os.File
			file, err = os.OpenFile(
				dbFullName,
				os.O_CREATE|os.O_APPEND|os.O_RDWR,
				0700,
			)
			if err != nil {
				return
			}
			_ = file.Close()
		}
	}()
	return errCh
}
