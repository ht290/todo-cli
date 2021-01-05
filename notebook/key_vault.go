package notebook

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"todo-cli/notebook/persistence"
)

const DefaultCredentialFilename = ".todo-config"

// Each notebook has a key for accessing its encrypted content.
// Keys can only be generated from users' plaintext passwords.
// Only the active user's key is stored in a cache, which will be cleared once the user returns the key.
type keyVault interface {
	// Create a new key for later use
	createKey(username, password string) error
	// Clear the key from the cache so no notebook is accessible.
	returnKey() error
	// Retrieve the key and store it in the cache, so the user can access its notebook.
	retrieveKey(username, password string) error
}

type ItemEncryptor interface {
	keyVault
	encrypt(plainItem Item) (encryptedItem Item, err error)
	decryptAll(encryptedItems []Item) (plainItems []Item, err error)
	decrypt(encryptedItems Item) (plainItem Item, err error)
	getActiveUser() UserId
}

type localKeyVault struct {
	activeUser     UserId
	encryptionKey  []byte
	credentialFile persistence.JsonFile
}

func UseLocalItemEncryptor(credentialFilename string) (*localKeyVault, error) {
	credentialFile := persistence.JsonFile{Filename: credentialFilename}
	var cred credential
	if err := credentialFile.Read(&cred); err != nil {
		return nil, err
	}
	return &localKeyVault{
		activeUser:     cred.UnlockedUsername,
		encryptionKey:  cred.UnlockedEncryptionKey,
		credentialFile: credentialFile,
	}, nil
}

type account struct {
	Username     string
	PasswordHash []byte
}

type credential struct {
	Accounts              []account
	UnlockedUsername      string
	UnlockedEncryptionKey []byte
}

func (p *localKeyVault) returnKey() error {
	var accountConfig credential
	if err := p.credentialFile.Read(&accountConfig); err != nil {
		return err
	}
	accountConfig.UnlockedEncryptionKey = nil
	accountConfig.UnlockedUsername = ""
	return p.credentialFile.Write(accountConfig)
}

func (p *localKeyVault) retrieveKey(username, password string) error {
	var accountConfig credential
	if err := p.credentialFile.Read(&accountConfig); err != nil {
		return err
	}
	// The salt can be any non-empty string, so that the notebook encryption key is different to password hash.
	const encryptionKeySalt = "salty"
	for _, account := range accountConfig.Accounts {
		if account.Username == username && comparePasswords(account.PasswordHash, []byte(password)) {
			keyHash, err := hashAndSalt([]byte(password + encryptionKeySalt))
			if err != nil {
				return err
			}
			accountConfig.UnlockedEncryptionKey = keyHash[0:32]
			accountConfig.UnlockedUsername = username
			p.activeUser = accountConfig.UnlockedUsername
			p.encryptionKey = accountConfig.UnlockedEncryptionKey
			return p.credentialFile.Write(accountConfig)
		}
	}
	return errors.New("invalid password")
}

func (p *localKeyVault) createKey(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}
	var accountConfig credential
	if err := p.credentialFile.Read(&accountConfig); err != nil {
		return err
	}
	for _, account := range accountConfig.Accounts {
		if account.Username == username {
			return errors.New("user already exists")
		}
	}
	passwordHash, err := hashAndSalt([]byte(password))
	if err != nil {
		return err
	}
	accountConfig.Accounts = append(accountConfig.Accounts, account{
		Username:     username,
		PasswordHash: passwordHash,
	})
	return p.credentialFile.Write(accountConfig)
}

func (p *localKeyVault) getActiveUser() UserId {
	return p.activeUser
}

func (p *localKeyVault) encrypt(item Item) (encryptedItem Item, err error) {
	item.EncryptedSummary, err = p.encryptData([]byte(item.Summary))
	if err != nil {
		return Item{}, err
	}
	item.Summary = ""
	return item, err
}

func (p *localKeyVault) decryptAll(items []Item) (plainItems []Item, err error) {
	for _, item := range items {
		decrypted, err := p.decrypt(item)
		if err != nil {
			return nil, err
		}
		plainItems = append(plainItems, decrypted)
	}
	return plainItems, err
}

func (p *localKeyVault) decrypt(item Item) (plainItem Item, err error) {
	summary, err := p.decryptData(item.EncryptedSummary)
	if err != nil {
		return Item{}, err
	}
	item.Summary = string(summary)
	item.EncryptedSummary = nil
	return item, nil
}

func comparePasswords(passwordHash []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(passwordHash, password)
	if err != nil {
		return false
	}
	return true
}

func hashAndSalt(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// Standard procedure for AES encryption
func (p *localKeyVault) encryptData(data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(p.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Standard procedure for AES decryption
func (p *localKeyVault) decryptData(data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(p.encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
