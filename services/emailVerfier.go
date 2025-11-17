package services

import (
	emailverifier "github.com/AfterShip/email-verifier"
)

func EmailVerification(email string) bool {

	verifier := emailverifier.NewVerifier()
	verifier = verifier.EnableDomainSuggest()
	verifier = verifier.EnableDomainSuggest()
	res, err := verifier.Verify(email)
	if res.Reachable == "no" {
		return false
	}
	if res.Suggestion != "" {
		return false
	}
	if !res.Syntax.Valid {
		return false
	}
	if err != nil {
		return false
	}
	if len(res.Email) == 0 {
		return false
	}
	if !res.HasMxRecords {
		return false
	}

	return true
}
