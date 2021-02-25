package zksneak_bn128

import (
	"ZKSneak/ZKSneak-crypto/bulletProofs/bp_bn128"
)

func Setup(b int64) (BulletProofSetupParams, error) {
	return bp_bn128.Setup(b)
}

func ProveTransfer(statement *ZKSneakTransferStatement, params *BulletProofSetupParams) (proof *ZKSneakTransferProof, err error) {
	proof = new(ZKSneakTransferProof)
	proof.ProveAnonEnc(statement.Relations)
	proof.ProveAnonRange(statement, params)
	proof.ProveEqual(statement.Relations)
	return proof, nil
}

func (proof *ZKSneakTransferProof) VerifyTransfer() bool {
	return proof.VerifyAnonEnc() && proof.VerifyAnonRange() && proof.VerifyEqual()
}