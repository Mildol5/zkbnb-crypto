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

package std

type RemoveLiquidityTx struct {
	/*
		- from account index
		- to account index
		- asset a id
		- asset a min amount
		- asset b id
		- asset b min amount
		- lp amount
		- gas account index
		- gas fee asset id
		- gas fee asset amount
	*/
	FromAccountIndex  uint32
	PairIndex         uint32
	AssetAId          uint32
	AssetAMinAmount   uint64
	AssetBId          uint32
	AssetBMinAmount   uint64
	LpAmount          uint64
	AssetAAmountDelta uint64
	AssetBAmountDelta uint64
	GasAccountIndex   uint32
	GasFeeAssetId     uint32
	GasFeeAssetAmount uint64
}

type RemoveLiquidityTxConstraints struct {
	FromAccountIndex  Variable
	PairIndex         Variable
	AssetAId          Variable
	AssetAMinAmount   Variable
	AssetBId          Variable
	AssetBMinAmount   Variable
	LpAmount          Variable
	AssetAAmountDelta Variable
	AssetBAmountDelta Variable
	GasAccountIndex   Variable
	GasFeeAssetId     Variable
	GasFeeAssetAmount Variable
}

func EmptyRemoveLiquidityTxWitness() (witness RemoveLiquidityTxConstraints) {
	return RemoveLiquidityTxConstraints{
		FromAccountIndex:  ZeroInt,
		PairIndex:         ZeroInt,
		AssetAId:          ZeroInt,
		AssetAMinAmount:   ZeroInt,
		AssetBId:          ZeroInt,
		AssetBMinAmount:   ZeroInt,
		LpAmount:          ZeroInt,
		AssetAAmountDelta: ZeroInt,
		AssetBAmountDelta: ZeroInt,
		GasAccountIndex:   ZeroInt,
		GasFeeAssetId:     ZeroInt,
		GasFeeAssetAmount: ZeroInt,
	}
}

func SetRemoveLiquidityTxWitness(tx *RemoveLiquidityTx) (witness RemoveLiquidityTxConstraints) {
	witness = RemoveLiquidityTxConstraints{
		FromAccountIndex:  tx.FromAccountIndex,
		PairIndex:         tx.PairIndex,
		AssetAId:          tx.AssetAId,
		AssetAMinAmount:   tx.AssetAMinAmount,
		AssetBId:          tx.AssetBId,
		AssetBMinAmount:   tx.AssetBMinAmount,
		LpAmount:          tx.LpAmount,
		AssetAAmountDelta: tx.AssetAAmountDelta,
		AssetBAmountDelta: tx.AssetBAmountDelta,
		GasAccountIndex:   tx.GasAccountIndex,
		GasFeeAssetId:     tx.GasFeeAssetId,
		GasFeeAssetAmount: tx.GasFeeAssetAmount,
	}
	return witness
}

func ComputeHashFromRemoveLiquidityTx(tx RemoveLiquidityTxConstraints, nonce Variable, hFunc MiMC) (hashVal Variable) {
	hFunc.Reset()
	hFunc.Write(
		tx.FromAccountIndex,
		tx.PairIndex,
		tx.AssetAId,
		tx.AssetAMinAmount,
		tx.AssetBId,
		tx.AssetBMinAmount,
		tx.LpAmount,
		tx.GasAccountIndex,
		tx.GasFeeAssetId,
		tx.GasFeeAssetAmount,
	)
	hFunc.Write(nonce)
	hashVal = hFunc.Sum()
	return hashVal
}

/*
	VerifyRemoveLiquidityTx:
	accounts order is:
	- FromAccount
		- Assets:
			- AssetA
			- AssetB
			- AssetGas
		- Liquidity:
			- LpAmount
	- ToAccount
		- Liquidity
			- AssetA
			- AssetB
	- GasAccount
		- Assets:
			- AssetGas
*/
func VerifyRemoveLiquidityTx(
	api API, flag Variable,
	tx *RemoveLiquidityTxConstraints,
	accountsBefore [NbAccountsPerTx]AccountConstraints, liquidityBefore LiquidityConstraints,
	hFunc *MiMC,
) {
	CollectPubDataFromRemoveLiquidity(api, flag, *tx, hFunc)
	// verify params
	// account index
	IsVariableEqual(api, flag, tx.FromAccountIndex, accountsBefore[0].AccountIndex)
	IsVariableEqual(api, flag, tx.GasAccountIndex, accountsBefore[1].AccountIndex)
	// asset id
	IsVariableEqual(api, flag, tx.AssetAId, accountsBefore[0].AssetsInfo[0].AssetId)
	IsVariableEqual(api, flag, tx.AssetBId, accountsBefore[0].AssetsInfo[1].AssetId)
	IsVariableEqual(api, flag, tx.GasFeeAssetId, accountsBefore[0].AssetsInfo[3].AssetId)
	IsVariableEqual(api, flag, tx.GasFeeAssetId, accountsBefore[2].AssetsInfo[0].AssetId)
	IsVariableEqual(api, flag, tx.AssetAId, liquidityBefore.AssetAId)
	IsVariableEqual(api, flag, tx.AssetBId, liquidityBefore.AssetBId)
	// should have enough lp
	IsVariableLessOrEqual(api, flag, tx.LpAmount, accountsBefore[0].AssetsInfo[2].LpAmount)
	// enough balance
	tx.AssetAMinAmount = UnpackAmount(api, tx.AssetAMinAmount)
	tx.AssetAAmountDelta = UnpackAmount(api, tx.AssetAAmountDelta)
	tx.AssetBMinAmount = UnpackAmount(api, tx.AssetBMinAmount)
	tx.AssetBAmountDelta = UnpackAmount(api, tx.AssetBAmountDelta)
	tx.LpAmount = UnpackAmount(api, tx.LpAmount)
	tx.GasFeeAssetAmount = UnpackFee(api, tx.GasFeeAssetAmount)
	IsVariableLessOrEqual(api, flag, tx.GasFeeAssetAmount, accountsBefore[0].AssetsInfo[0].Balance)
	// TODO verify LP
	Delta_LPCheck := api.Mul(tx.AssetAAmountDelta, tx.AssetBAmountDelta)
	LPCheck := api.Mul(tx.LpAmount, tx.LpAmount)
	IsVariableLessOrEqual(api, flag, Delta_LPCheck, LPCheck)
	IsVariableLessOrEqual(api, flag, tx.AssetAMinAmount, tx.AssetAAmountDelta)
	IsVariableLessOrEqual(api, flag, tx.AssetBMinAmount, tx.AssetBAmountDelta)
}