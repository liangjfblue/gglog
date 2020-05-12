/**
 *
 * @author liangjf
 * @create on 2020/5/12
 * @version 1.0
 */
package kafka

import "log"

type DefaultKafkaCallback struct {
}

func (callback *DefaultKafkaCallback) Success(msg string) {
	log.Printf("SuccessCallback: %s\n", msg)
}

func (callback *DefaultKafkaCallback) Fail(err error) {
	log.Printf("FailCallback: %s\n", err.Error())
}
