/*

  Copyright 2019 Dmitry Kolesnikov, All Rights Reserved

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package mock

import (
	"fmt"
	µ "github.com/fogfish/gouldian/v2"
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
