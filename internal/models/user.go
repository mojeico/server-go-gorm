package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type User struct {
	gorm.Model
	CompanyID  int           `json:"company_id"`
	Name       string        `json:"name"`
	Email      string        `json:"email" binding:"required"`
	Phone      string        `json:"phone"`
	Username   string        `json:"username" binding:"required"`
	Password   string        `json:"password" binding:"required"`
	Groups     pq.Int32Array `json:"groups" gorm:"type:integer[]"`
	AuthStatus bool          `json:"auth_status"`
	Status     string        `json:"status"` // activated - confirm - deleted  -- login only for confirm
	IsActive   bool          `json:"is_active"`
	IsDeleted  bool          `json:"is_deleted"`

	PictureName string `json:"picture_name"`

	Role string `json:"role" gorm:"default=USER"`

	LastConnection time.Time `json:"last_connection"`
	UserAgent      string    `json:"user_agent"`
	IpAddress      string    `json:"ip_address"`
}

type UpdateUserInput struct {
	CompanyID  int           `json:"company_id"`
	Name       string        `json:"name"`
	Email      string        `json:"email"`
	Phone      string        `json:"phone"`
	Username   string        `json:"username"`
	Password   string        `json:"password"`
	Groups     pq.Int32Array `json:"groups" gorm:"type:integer[]"`
	AuthStatus bool          `json:"auth_status"`
	IsActive   bool          `json:"is_active"`

	PictureName string `json:"picture_name"`

	LastConnection time.Time `json:"last_connection"`
	UserAgent      string    `json:"user_agent"`
	IpAddress      string    `json:"ip_address"`
}

const (
	minRequiredLen = 4
	maxRequiredLen = 32
)

type UserIpUpdate struct {
	LastConnection time.Time `json:"last_connection"`
	UserAgent      string    `json:"user_agent"`
	IpAddress      string    `json:"ip_address"`
}

var regexpEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (u User) IsValid() error {
	if u.Username == "" {
		return errors.New("'username' is required")
	}

	if len(u.Username) < minRequiredLen || len(u.Username) > maxRequiredLen {
		return fmt.Errorf("the name field must be between %d-%d chars", minRequiredLen, maxRequiredLen)
	}

	if u.Email == "" {
		return errors.New("the email field is required")
	}

	if !regexpEmail.Match([]byte(u.Email)) {
		return errors.New("the email field should be a valid email address")
	}

	if u.Password == "" {
		return errors.New("'password' is required")
	}

	if len(u.Password) < minRequiredLen*2 || len(u.Password) > maxRequiredLen {
		return fmt.Errorf("the password field must be between %d-%d chars", minRequiredLen, maxRequiredLen)
	}

	return nil
}
