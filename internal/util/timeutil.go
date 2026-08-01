package util

import "time"

// BeijingLocation 北京时间（UTC+8）时区。使用固定时区偏移，避免依赖操作系统的时区数据库，
// 保证在本地（如 Windows）和服务器上表现一致。
var BeijingLocation = time.FixedZone("CST", 8*3600)

// timeLayout 统一的时间字符串格式。
const timeLayout = "2006-01-02 15:04:05"

// ParseBeijingTime 将 "2006-01-02 15:04:05" 格式的北京时间字符串解析为 time.Time。
// 注意：直接用 time.Parse 会把无时区的字符串当作 UTC，导致时刻偏差 8 小时，
// 因此这里必须用 ParseInLocation 明确按北京时区解析。
func ParseBeijingTime(s string) (time.Time, error) {
	return time.ParseInLocation(timeLayout, s, BeijingLocation)
}

// FormatBeijingTime 将时间转换为北京时间后按 "2006-01-02 15:04:05" 格式化。
func FormatBeijingTime(t time.Time) string {
	return t.In(BeijingLocation).Format(timeLayout)
}
