package utils

func CheckCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	for _, c := range cep {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}
