/**
 * @Author: Administrator
 * @Description:
 * @File:  redis_test
 * @Version: 1.0.0
 * @Date: 2019/12/11 17:32
 */

package pkg

import "testing"

func TestIncrBy(t *testing.T) {
	IncrBy("test", 3)
}
