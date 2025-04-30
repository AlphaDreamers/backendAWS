package delivery

import (
	"github.com/SwanHtetAungPhyo/auth/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Delivery interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
}

var _Delivery = (*UserDelivery)(nil)

type UserDelivery struct {
	log     *logrus.Logger
	service service.Service
}

func NewUserDelivery(log *logrus.Logger, service2 service.Service) *UserDelivery {
	return &UserDelivery{
		log:     log,
		service: service2,
	}
}

func (u UserDelivery) Login(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserDelivery) Register(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserDelivery) Logout(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserDelivery) Refresh(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserDelivery) ForgotPassword(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserDelivery) ResetPassword(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserDelivery) GetProfile(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
