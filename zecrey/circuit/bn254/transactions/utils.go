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

package transactions

import (
	"github.com/zecrey-labs/zecrey-crypto/zecrey/circuit/bn254/std"
)

func SelectCommonPart(
	api API,
	flag Variable,
	c, cCheck Variable,
	pkProofs, pkProofsCheck [MaxRangeProofCount]std.CommonPkProof,
	tProofs, tProofsCheck [MaxRangeProofCount]std.CommonTProof,
) (cRes Variable, pkProofsRes [MaxRangeProofCount]std.CommonPkProof, tProofsRes [MaxRangeProofCount]std.CommonTProof) {
	cRes = api.Select(flag, c, cCheck)
	for i := 0; i < MaxRangeProofCount; i++ {
		pkProofsRes[i] = std.SelectCommonPkProof(api, flag, pkProofs[i], pkProofsCheck[i])
		tProofsRes[i] = std.SelectCommonTProof(api, flag, tProofs[i], tProofsCheck[i])
	}
	return cRes, pkProofsRes, tProofsRes
}

func SelectDeltas(
	api API,
	flag Variable,
	deltas, deltasCheck [NbAccountsPerTx]std.AccountDeltaConstraints,
) (deltasRes [NbAccountsPerTx]std.AccountDeltaConstraints) {
	for i := 0; i < NbAccountsPerTx; i++ {
		for j := 0; j < NbAccountAssetsPerAccount; j++ {
			deltasRes[i].AssetsDeltaInfo[j] =
				std.SelectElgamal(api, flag, deltas[i].AssetsDeltaInfo[j], deltasCheck[i].AssetsDeltaInfo[j])
		}
		deltasRes[i].LockedAssetDeltaInfo = api.Select(flag, deltas[i].LockedAssetDeltaInfo, deltasCheck[i].LockedAssetDeltaInfo)
		deltasRes[i].LiquidityDeltaInfo.AssetADelta =
			api.Select(flag, deltas[i].LiquidityDeltaInfo.AssetADelta, deltasCheck[i].LiquidityDeltaInfo.AssetADelta)
		deltasRes[i].LiquidityDeltaInfo.AssetARDelta =
			api.Select(flag, deltas[i].LiquidityDeltaInfo.AssetARDelta, deltasCheck[i].LiquidityDeltaInfo.AssetARDelta)
		deltasRes[i].LiquidityDeltaInfo.AssetBDelta =
			api.Select(flag, deltas[i].LiquidityDeltaInfo.AssetBDelta, deltasCheck[i].LiquidityDeltaInfo.AssetBDelta)
		deltasRes[i].LiquidityDeltaInfo.AssetBRDelta =
			api.Select(flag, deltas[i].LiquidityDeltaInfo.AssetBRDelta, deltasCheck[i].LiquidityDeltaInfo.AssetBRDelta)
		deltasRes[i].LiquidityDeltaInfo.LpEncDelta =
			std.SelectElgamal(api, flag, deltas[i].LiquidityDeltaInfo.LpEncDelta, deltasCheck[i].LiquidityDeltaInfo.LpEncDelta)
	}
	return deltasRes
}

func GetAccountDeltasFromDepositTx(
	api API, tool *EccTool, h Point,
	proof DepositOrLockTxConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	CDelta := tool.ZeroElgamalEnc()
	CDelta.CR = tool.ScalarMul(h, proof.Amount)

	// account A
	deltas[DepositFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			CDelta,
			CDelta,
			CDelta,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[1] = deltas[0]
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromLockTx(
	api API, tool *EccTool,
	proof DepositOrLockTxConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// account A
	deltas[TransferAccountA] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
		},
		// locked asset
		LockedAssetDeltaInfo: proof.Amount,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[1] = deltas[0]
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromUnlockProof(
	api API, tool *EccTool,
	proof UnlockProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[UnlockFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_Delta,
			// gas asset
			proof.C_fee_DeltaForFrom,
			proof.C_Delta,
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(proof.DeltaAmount),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[UnlockGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromTransferProof(
	api API, tool *EccTool,
	proof TransferProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// account A
	deltas[TransferAccountA] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.SubProofs[TransferAccountA].CDelta,
			proof.SubProofs[TransferAccountA].CDelta,
			proof.SubProofs[TransferAccountA].CDelta,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// account B
	deltas[TransferAccountB] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.SubProofs[TransferAccountB].CDelta,
			proof.SubProofs[TransferAccountB].CDelta,
			proof.SubProofs[TransferAccountB].CDelta,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// account C
	deltas[TransferAccountC] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.SubProofs[TransferAccountC].CDelta,
			proof.SubProofs[TransferAccountC].CDelta,
			proof.SubProofs[TransferAccountC].CDelta,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[TransferGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	return deltas
}

func GetAccountDeltasFromSwapProof(
	api API, tool *EccTool,
	proof SwapProofConstraints, poolAccount std.AccountConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	deltas[SwapFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_uA_Delta,
			// to asset
			proof.C_uB_Delta,
			// gas asset
			proof.C_fee_DeltaForFrom,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// pool account
	B_A_Delta := api.Sub(proof.B_A_Delta, proof.B_treasuryfee_Delta)
	B_B_Delta := proof.B_B_Delta
	A_R := proof.R_DeltaA
	B_R := proof.R_DeltaB
	isSameAsset := api.IsZero(api.Sub(proof.AssetAId, poolAccount.LiquidityInfo.AssetAId))
	Pool_A_Delta := api.Select(isSameAsset, B_A_Delta, api.Neg(B_B_Delta))
	Pool_B_Delta := api.Select(isSameAsset, api.Neg(B_B_Delta), B_A_Delta)
	Pool_A_R := api.Select(isSameAsset, A_R, B_R)
	Pool_B_R := api.Select(isSameAsset, B_R, A_R)
	deltas[SwapPoolAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		// pool info
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  Pool_A_Delta,
			AssetBDelta:  Pool_B_Delta,
			AssetARDelta: Pool_A_R,
			AssetBRDelta: Pool_B_R,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// treasury account
	deltas[SwapTreasuryAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_treasuryfee_Delta,
			proof.C_treasuryfee_Delta,
			proof.C_treasuryfee_Delta,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[SwapGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	return deltas
}

func GetAccountDeltasFromAddLiquidityProof(
	api API, tool *EccTool,
	proof AddLiquidityProofConstraints, poolAccount std.AccountConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	deltas[AddLiquidityFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_uA_Delta,
			// to asset
			proof.C_uB_Delta,
			// gas asset
			proof.C_fee_DeltaForFrom,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   proof.C_LP_Delta,
		},
	}
	// pool account
	B_A_Delta := proof.B_A_Delta
	B_B_Delta := proof.B_B_Delta
	isSameAsset := api.IsZero(api.Sub(proof.AssetAId, poolAccount.LiquidityInfo.AssetAId))
	Pool_A_Delta := api.Select(isSameAsset, B_A_Delta, B_B_Delta)
	Pool_B_Delta := api.Select(isSameAsset, B_B_Delta, B_A_Delta)
	Pool_A_R := api.Select(isSameAsset, proof.R_DeltaA, proof.R_DeltaB)
	Pool_B_R := api.Select(isSameAsset, proof.R_DeltaB, proof.R_DeltaA)
	deltas[AddLiquidityPoolAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		// pool info
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  Pool_A_Delta,
			AssetBDelta:  Pool_B_Delta,
			AssetARDelta: Pool_A_R,
			AssetBRDelta: Pool_B_R,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[AddLiquidityGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromRemoveLiquidityProof(
	api API, tool *EccTool,
	proof RemoveLiquidityProofConstraints, poolAccount std.AccountConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	deltas[RemoveLiquidityFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_uA_Delta,
			// to asset
			proof.C_uB_Delta,
			// gas asset
			proof.C_fee_DeltaForFrom,
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   proof.C_u_LP_Delta,
		},
	}
	// pool account
	B_A_Delta := api.Neg(proof.B_A_Delta)
	B_B_Delta := api.Neg(proof.B_B_Delta)
	isSameAsset := api.IsZero(api.Sub(proof.AssetAId, poolAccount.LiquidityInfo.AssetAId))
	Pool_A_Delta := api.Select(isSameAsset, B_A_Delta, B_B_Delta)
	Pool_B_Delta := api.Select(isSameAsset, B_B_Delta, B_A_Delta)
	Pool_A_R := api.Select(isSameAsset, proof.R_DeltaA, proof.R_DeltaB)
	Pool_B_R := api.Select(isSameAsset, proof.R_DeltaB, proof.R_DeltaA)
	deltas[RemoveLiquidityPoolAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
		},
		// locked asset
		LockedAssetDeltaInfo: std.ZeroInt,
		// pool info
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  Pool_A_Delta,
			AssetBDelta:  Pool_B_Delta,
			AssetARDelta: Pool_A_R,
			AssetBRDelta: Pool_B_R,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[RemoveLiquidityGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromWithdrawProof(
	api API, tool *EccTool,
	proof WithdrawProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[WithdrawFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// from asset
			proof.C_Delta,
			// gas asset
			proof.C_fee_DeltaForFrom,
			proof.C_Delta,
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(std.ZeroInt),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[WithdrawGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromDepositNftProof(
	api API, tool *EccTool,
	proof DepositNftTxConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[0] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(std.ZeroInt),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[1] = deltas[0]
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromMintNftProof(
	api API, tool *EccTool,
	proof ClaimNftProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[MintNftFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(std.ZeroInt),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[MintNftGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromTransferNftProof(
	api API, tool *EccTool,
	proof ClaimNftProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[MintNftFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(std.ZeroInt),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[MintNftGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[MintNftToAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromSetNftPriceProof(
	api API, tool *EccTool,
	proof SetNftPriceProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[SetNftPriceFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(std.ZeroInt),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[SetNftPriceGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromBuyNftProof(
	api API, tool *EccTool,
	proof BuyNftProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[BuyNftFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_Delta,
			proof.C_fee_DeltaForFrom,
			proof.C_Delta,
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(std.ZeroInt),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[BuyNftGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[BuyNftToAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
			tool.ZeroElgamalEnc(),
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[3] = deltas[0]
	return deltas
}

func GetAccountDeltasFromWithdrawNftProof(
	api API, tool *EccTool,
	proof WithdrawNftProofConstraints,
) (deltas [NbAccountsPerTx]std.AccountDeltaConstraints) {
	// from account
	// from asset
	deltas[WithdrawNftFromAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
			proof.C_fee_DeltaForFrom,
		},
		// locked asset
		LockedAssetDeltaInfo: api.Neg(std.ZeroInt),
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	// gas account
	deltas[WithdrawNftGasAccount] = std.AccountDeltaConstraints{
		AssetsDeltaInfo: [3]ElGamalEncConstraints{
			// gas asset
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
			proof.C_fee_DeltaForGas,
		},
		LockedAssetDeltaInfo: std.ZeroInt,
		LiquidityDeltaInfo: std.AccountLiquidityDeltaConstraints{
			AssetADelta:  std.ZeroInt,
			AssetBDelta:  std.ZeroInt,
			AssetARDelta: std.ZeroInt,
			AssetBRDelta: std.ZeroInt,
			LpEncDelta:   tool.ZeroElgamalEnc(),
		},
	}
	deltas[2] = deltas[0]
	deltas[3] = deltas[0]
	return deltas
}