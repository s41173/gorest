package services

import (
	"fmt"
	"go-rest/config"
	"go-rest/models"
	"go-rest/utils"
)

type CustomerService struct{}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) ChangePassword(username, newpassword string) (int, error) { // return otp, userid, err

	customer := models.Customer{}

	var res *models.Customer
	validUser := customer.CheckUser(config.DB, username)
	validPhone := customer.CheckUserPhone(config.DB, username)

	if !validUser && !validPhone {
		return 400, fmt.Errorf("User Not Found")
	}

	if validUser {
		res = customer.GetByUsername(config.DB, username)
	} else if validPhone {
		res = customer.GetByPhone(config.DB, username)
	}

	if utils.VerifyPassword(newpassword, res.Password) == true {
		return 403, fmt.Errorf("Can't use previous password...!")
	}

	err := models.UpdatePassword(config.DB, int(res.ID), newpassword)
	if err != nil {
		return 500, fmt.Errorf("Failed to update password")
	}

	return 200, nil
}
