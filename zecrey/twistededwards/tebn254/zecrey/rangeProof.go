/*
 * Copyright © 2021 Zecrey Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package zecrey

import (
	"log"
	"math/big"
	"zecrey-crypto/rangeProofs/twistededwards/tebn254/ctrange"
)

func proveCtRangeRoutine(b int64, g, h *Point, r *big.Int, proof *RangeProof, swapRangeChan chan int) {
	var (
		err error
	)
	bar_r, rangeProof, err := ctrange.Prove(b, g, h)
	if err != nil {
		log.Println("[proveCtRangeRoutine] err info:", err)
		swapRangeChan <- ErrCode
		return
	}
	*proof = *rangeProof
	*r = *bar_r
	swapRangeChan <- 1
}
