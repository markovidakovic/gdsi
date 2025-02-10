package sec

import "golang.org/x/crypto/bcrypt"

func EncryptPwd(pwd string) (string, error) {
	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwdBytes), nil
}
