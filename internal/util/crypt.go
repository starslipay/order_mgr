package util

import (
	"crypto/md5"
	"fmt"
	"strconv"

	"github.com/starslipay/order_mgr/internal/xerr"
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"
)

func GenMD5(input string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(input)))
}

func GenOrderSuccessToken(transaction_id, out_order_no string, merchant_uid, uid, amount int64) string {
	md5Str := GenMD5(transaction_id + "|" + out_order_no + "|" + strconv.FormatInt(merchant_uid, 10) +
		"|" + strconv.FormatInt(uid, 10) + "|" + strconv.FormatInt(amount, 10))
	return md5Str
}

func GenC2BDeductToken(transaction_id string, merchant_uid, uid, amount int64) string {
	md5Str := GenMD5(transaction_id + "|" + strconv.FormatInt(merchant_uid, 10) +
		"|" + strconv.FormatInt(uid, 10) + "|" + strconv.FormatInt(amount, 10))
	return md5Str
}

func CheckC2BDeductToken(transaction_id string, merchant_uid, uid, amount int64, token string) error {
	if GenC2BDeductToken(transaction_id, merchant_uid, uid, amount) != token {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeCheckDeductToken, "deductToken check failed")
	}
	return nil
}
