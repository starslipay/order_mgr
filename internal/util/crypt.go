package util

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func GenMD5(input string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(input)))
}

func GenOrderSuccessToken(transaction_id, spid, user_id string, uid, amount int64) string {
	md5Str := GenMD5(transaction_id + "|" + spid + "|" + user_id + "|" + strconv.FormatInt(uid, 10) + "|" + strconv.FormatInt(amount, 10))
	return md5Str
}
