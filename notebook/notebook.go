package notebook

type Notebook interface {
	Items
	Create(username, password string) error
	Unlock(username, password string) error
	Lock() error
}

func New(encryptor ItemEncryptor, storage ItemStorage) (*multiuserNotebook, error) {
	return &multiuserNotebook{
		storage:   storage,
		encryptor: encryptor,
	}, nil
}
