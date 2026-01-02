package entity

type Customer struct {
	Id          string
	Name        string
	Email       string
	CPF         string
	IsAnonymous bool
}

func (c Customer) Build() Customer {
	return Customer{
		Id:          c.Id,
		Name:        c.Name,
		Email:       c.Email,
		CPF:         c.CPF,
		IsAnonymous: c.IsAnonymous,
	}
}
