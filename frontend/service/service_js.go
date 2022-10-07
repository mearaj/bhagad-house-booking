//go:build js

package service

import (
	"github.com/mearaj/bhagad-house-booking/common/alog"
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
	}()
	return errCh
}
