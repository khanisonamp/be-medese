package notifyline

import (
	"fmt"

	"github.com/juunini/simple-go-line-notify/notify"
	"github.com/sirupsen/logrus"
)

func SendMsgErrToLine(dateStart, dateEnd, productCode, stock, quantity, quantityManual, remainingToday string) error {
	if err := notify.SendText("ilecOF6w7iKNRA9r664BKXZu8ZX0CnENTBwTWDU41os", fmt.Sprintf("ข้อมูลวันที่ :\n%s ถึง %s\nรหัสสินค้า : %s\nสต๊อกคงเหลือ : %s\nออเดอร์ : %s\nออเดอร์แมนนวล : %s\nคงเหลือล่าสุด : %s", dateStart, dateEnd, productCode, stock, quantity, quantityManual, remainingToday)); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
