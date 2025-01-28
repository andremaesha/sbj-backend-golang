package tests

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	Save(user User) error
	FindByID(id int) (*User, error)
}

type UserService struct {
	Repo UserRepository
}

func (u *UserService) RegisterUser(user User) error {
	return u.Repo.Save(user)
}
