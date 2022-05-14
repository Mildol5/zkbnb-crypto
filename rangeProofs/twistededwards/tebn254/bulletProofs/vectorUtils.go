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

package bulletProofs

import (
	"errors"
	"math/big"
	curve "github.com/zecrey-labs/zecrey-crypto/ecc/ztwistededwards/tebn254"
	"github.com/zecrey-labs/zecrey-crypto/ffmath"
)

/*
SampleRandomVector generates a vector composed by random big numbers.
*/
func RandomVector(N int64) []*big.Int {
	s := make([]*big.Int, N)
	for i := int64(0); i < N; i++ {
		s[i] = curve.RandomValue()
	}
	return s
}

/*
VectorCopy returns a vector composed by copies of a.
*/
func VectorCopy(a *big.Int, n int64) ([]*big.Int, error) {
	var (
		i      int64
		result []*big.Int
	)
	result = make([]*big.Int, n)
	i = 0
	for i < n {
		result[i] = a
		i = i + 1
	}
	return result, nil
}

/*
VectorConvertToBig converts an array of int64 to an array of big.Int.
*/
func ToBigIntVec(a []int64, n int64) ([]*big.Int, error) {
	var (
		i      int64
		result []*big.Int
	)
	result = make([]*big.Int, n)
	i = 0
	for i < n {
		result[i] = big.NewInt(a[i])
		i = i + 1
	}
	return result, nil
}

/*
VectorSub computes vector addition componentwisely.
*/
func VectorSub(a, b []*big.Int) ([]*big.Int, error) {
	var (
		result  []*big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("size of first argument is different from size of second argument")
	}
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = ffmath.SubMod(a[i], b[i], Order)
		i = i + 1
	}
	return result, nil
}

/*
VectorMul computes vector multiplication componentwisely.
*/
func VectorMul(a, b []*big.Int) ([]*big.Int, error) {
	var (
		result  []*big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("size of first argument is different from size of second argument")
	}
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = ffmath.MultiplyMod(a[i], b[i], Order)
		i = i + 1
	}
	return result, nil
}

/*
ScalarProduct return the inner product between a and b.
*/
func ScalarVecMul(a, b []*big.Int) (*big.Int, error) {
	var (
		result  *big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("size of first argument is different from size of second argument")
	}
	i = 0
	result = big.NewInt(0)
	for i < n {
		ab := ffmath.MultiplyMod(a[i], b[i], Order)
		result = ffmath.AddMod(result, ab, Order)
		i = i + 1
	}
	return result, nil
}

/*
VectorAdd computes vector addition componentwisely.
*/
func VectorAdd(a, b []*big.Int) ([]*big.Int, error) {
	var (
		result  []*big.Int
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("size of first argument is different from size of second argument")
	}
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = ffmath.AddMod(a[i], b[i], Order)
		i = i + 1
	}
	return result, nil
}

/*
VectorScalarMul computes vector scalar multiplication componentwisely.
*/
func VectorScalarMul(a []*big.Int, b *big.Int) ([]*big.Int, error) {
	var (
		result []*big.Int
		i, n   int64
	)
	n = int64(len(a))
	i = 0
	result = make([]*big.Int, n)
	for i < n {
		result[i] = ffmath.MultiplyMod(a[i], b, Order)
		i = i + 1
	}
	return result, nil
}

/*
VectorECMul computes vector EC addition componentwisely.
*/
func VectorECAdd(a, b []*Point) ([]*Point, error) {
	var (
		result  []*Point
		i, n, m int64
	)
	n = int64(len(a))
	m = int64(len(b))
	if n != m {
		return nil, errors.New("size of first argument is different from size of second argument")
	}
	result = make([]*Point, n)
	i = 0
	for i < n {
		result[i] = curve.Add(a[i], b[i])
		i = i + 1
	}
	return result, nil
}

/*
VectorExp computes Prod_i^n{a[i]^b[i]}.
*/
func VectorExp(a []*Point, b []*big.Int) (result *Point, err error) {
	n := int64(len(a))
	m := int64(len(b))
	if n < m {
		return nil, errors.New("size of first argument is different from size of second argument")
	}
	i := int64(0)
	res := curve.ZeroPoint()
	for i < m {
		res.Add(res, curve.ScalarMul(a[i], b[i]))
		i = i + 1
	}
	return res, nil
}

/*
VectorScalarExp computes a[i]^b for each i.
*/
func vectorScalarExp(a []*Point, b *big.Int) []*Point {
	var (
		result []*Point
		n      int64
	)
	n = int64(len(a))
	result = make([]*Point, n)
	for i := int64(0); i < n; i++ {
		result[i] = curve.ScalarMul(a[i], b)
	}
	return result
}