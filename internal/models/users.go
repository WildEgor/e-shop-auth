package models

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ActiveStatus  = "active"
	DeletedStatus = "deleted"
	BlockedStatus = "blocked"
)

type UsersModel struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Email        string             `json:"email,omitempty" bson:"email"`
	Phone        string             `json:"phone,omitempty" bson:"phone"`
	Password     string             `json:"password,omitempty" bson:"password"`
	Verification VerificationModel  `json:"verification,omitempty" bson:"verification"`
	OTP          OTPModel           `json:"otp,omitempty" bson:"otp"`
	Status       string             `json:"status,omitempty" bson:"status"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt    time.Time          `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type VerificationModel struct {
	NewPhone     string    `json:"new_phone" bson:"new_phone"`
	NewPhoneCode string    `json:"new_phone_code" bson:"new_phone_code"`
	NewPhoneDate time.Time `json:"new_phone_date,omitempty" bson:"new_phone_date"`
	NewEmail     string    `json:"new_email" bson:"new_email"`
	NewEmailCode string    `json:"new_email_code" bson:"new_email_code"`
	NewEmailDate time.Time `json:"new_email_date,omitempty" bson:"new_email_date"`
}

type OTPModel struct {
	Identity string    `json:"identity" bson:"identity"`
	Code     string    `json:"code" bson:"code"`
	ExpireAt time.Time `json:"expire_at" bson:"expire_at"`
}

func (us *UsersModel) ComparePassword(password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(us.Password), []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}

func (us *UsersModel) SetInfo(firstname, lastname string) {
	us.FirstName = firstname
	us.LastName = lastname
}

func (us *UsersModel) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}

	us.Password = string(hash)

	return nil
}

func (us *UsersModel) IsOTPResendAvailable() error {
	if time.Since(us.OTP.ExpireAt) < 2*time.Minute {
		return errors.New("too many otp")
	}

	return nil
}

func (us *UsersModel) IsEmailConfirmResendAvailable() error {
	if time.Since(us.Verification.NewEmailDate) < 2*time.Minute {
		return errors.New("too many email confirm")
	}

	return nil
}

func (us *UsersModel) IsPhoneConfirmResendAvailable() error {
	if time.Since(us.Verification.NewPhoneDate) < 2*time.Minute {
		return errors.New("too many phone confirm")
	}

	return nil
}

func (us *UsersModel) VerifyOTP(phone string, code string) error {
	if us.OTP.Identity != phone {
		return errors.New("empty identity")
	}

	if us.OTP.Code != code {
		return errors.New("wrong code")
	}

	if us.OTP.ExpireAt.Before(time.Now()) {
		return errors.New("code expired")
	}

	return nil
}

func (us *UsersModel) ClearOTP() {
	us.OTP = OTPModel{}
}

func (us *UsersModel) VerifyIdentity(identity, code string) error {
	if us.Verification.NewEmail == identity {
		if us.Verification.NewEmailCode != code {
			return errors.New("wrong code")
		}

		if us.Verification.NewEmailDate.Before(time.Now()) {
			return errors.New("code expired")
		}

		return nil
	}

	if us.Verification.NewPhone == identity {
		if us.Verification.NewPhone != code {
			return errors.New("wrong code")
		}

		if us.Verification.NewPhoneDate.Before(time.Now()) {
			return errors.New("code expired")
		}

		return nil
	}

	return errors.New("wrong code")
}

func (us *UsersModel) ClearEmailVerification() {
	us.Verification.NewEmail = ""
	us.Verification.NewEmailCode = ""
	us.Verification.NewEmailDate = time.Now()
}

func (us *UsersModel) ClearPhoneVerification() {
	us.Verification.NewPhone = ""
	us.Verification.NewPhoneCode = ""
	us.Verification.NewPhoneDate = time.Now()
}

func (us *UsersModel) IsEmailEqual(identity string) bool {
	return us.Email == identity
}

func (us *UsersModel) IsPhoneEqual(identity string) bool {
	return us.Phone == identity
}

func (us *UsersModel) UpdateOTP(identity, code string) {
	us.OTP.Identity = identity
	us.OTP.Code = code
	us.OTP.ExpireAt = time.Now().Add(time.Minute * 5) // TODO: make configurable
}

func (us *UsersModel) UpdatePhone(phone string) {
	us.Phone = phone
}

func (us *UsersModel) UpdateEmail(email string) {
	us.Email = email
}

func (us *UsersModel) UpdateEmailVerification(email, code string) {
	us.Verification.NewEmail = email
	us.Verification.NewEmailCode = code
	us.Verification.NewEmailDate = time.Now().Add(time.Minute * 5) // TODO: make configurable
}

func (us *UsersModel) UpdatePhoneVerification(phone, code string) {
	us.Verification.NewPhone = phone
	us.Verification.NewPhoneCode = code
	us.Verification.NewPhoneDate = time.Now().Add(time.Minute * 5) // TODO: make configurable
}

func (us *UsersModel) IsActive() bool {
	return us.Status == ActiveStatus
}
