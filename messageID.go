package tcpx

import "strconv"

const (
	SERVER_ERROR = 500
	CLIENT_ERROR = 400
	OK           = 200
	NOT_AUTH     = 403
)

func MessageWrap(messageId int32, srvCode int32) int32 {
	rs := strconv.FormatInt(int64(srvCode), 10) + strconv.FormatInt(int64(messageId), 10)

	r, _ := strconv.ParseInt(rs, 10, 64)
	return int32(r)
}
