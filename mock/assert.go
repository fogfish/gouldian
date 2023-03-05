package mock

import (
	"fmt"
	µ "github.com/fogfish/gouldian"
)

//
// assert functions for it/v2
//

func CheckStatusCode(err error, code int) error {
	switch v := err.(type) {
	case *µ.Output:
		if v.Status == code {
			return nil
		}
		return fmt.Errorf("status code %v be equal to %v", v.Status, code)

	default:
		return fmt.Errorf("type of %T be *µ.Output", err)
	}
}

func CheckOutput(err error, expect string) error {
	assert := fmt.Errorf("%v be equal to %v", err.Error(), expect)

	if err.Error() != expect {
		return assert
	}
	return nil
}
