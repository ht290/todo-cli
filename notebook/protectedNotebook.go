package notebook

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"todo-cli/notebook/internal/persistence"
)

const defaultAccountFilename = ".todo-config"

type protectedNotebook struct {
	username      string
	encryptionKey string
	accountFile   persistence.JsonFile
}

func newProtectedNotebook(accountFilename string) protectedNotebook {
	return protectedNotebook{accountFile: persistence.JsonFile{Filename: accountFilename}}
}

type account struct {
	Username     string
	PasswordHash string
	// The EncryptionKey will be erased in the config file if the account is locked
	EncryptionKey string
}

func (p *protectedNotebook) Lock() error {
	var accounts []account
	if err := p.accountFile.Read(&accounts); err != nil {
		return err
	}
	for i := range accounts {
		accounts[i].EncryptionKey = ""
	}
	return p.accountFile.Write(accounts)
}

func (p *protectedNotebook) Unlock(username, password string) error {
	var accounts []account
	if err := p.accountFile.Read(&accounts); err != nil {
		return err
	}
	// The salt can be any non-empty string, so that the notebook encryption key is different to password hash.
	const encryptionKeySalt = "salty"
	for i, account := range accounts {
		if account.Username == username && comparePasswords(account.PasswordHash, password) {
			var err error
			accounts[i].EncryptionKey, err = hashAndSalt(password + encryptionKeySalt)
			if err != nil {
				return err
			}
			return p.accountFile.Write(accounts)
		}
	}
	return errors.New("invalid password")
}

func (p *protectedNotebook) Create(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}
	var accounts []account
	if err := p.accountFile.Read(&accounts); err != nil {
		return err
	}
	for _, account := range accounts {
		if account.Username == username {
			return errors.New("user already exists")
		}
	}
	passwordHash, err := hashAndSalt(password)
	if err != nil {
		return err
	}
	accounts = append(accounts, account{
		Username:     username,
		PasswordHash: passwordHash,
	})
	return p.accountFile.Write(accounts)
}

func comparePasswords(passwordHash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
