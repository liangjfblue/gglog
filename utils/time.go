/**
 *
 * @author liangjf
 * @create on 2020/5/12
 * @version 1.0
 */
package utils

import "time"

// GMT location
var gmtLoc = time.FixedZone("GMT", 0)

// NowRFC1123 returns now time in RFC1123 format with GMT timezone,
// eg. "Mon, 02 Jan 2006 15:04:05 GMT".
func NowRFC1123() string {
	return time.Now().In(gmtLoc).Format(time.RFC1123)
}
