package domains

import (
	core_dtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	"github.com/gofiber/fiber/v3"
)

const defErrCode = 99

var errCodesMessages = map[int]string{
	99:  "unknown error",
	100: "empty request",
	101: "error request read",
	102: "error request unmarshal",
	103: "incorrect data",
	104: "incorrect method",
	105: "token not found",
	106: "user not found",
	113: "email already exists",
	114: "phone already exists",
	115: "resend not available",
	116: "incorrect code",
}

func SetInternalServerStatus(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusInternalServerError)
	r.SetError(99, errCodesMessages[99])
}

func SetMalformedCodeError(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusBadRequest)
	r.SetError(116, errCodesMessages[116])
}

func SetSendCodeError(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusInternalServerError)
	r.SetError(113, errCodesMessages[113]) // TODO: change
}

func SetSendCodeTimeoutError(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusBadRequest)
	r.SetError(115, errCodesMessages[115])
}

func SetEmailAlreadyExistError(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusBadRequest)
	r.SetError(113, errCodesMessages[113]) // TODO: change
}

func SetPhoneAlreadyExistError(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusBadRequest)
	r.SetError(114, errCodesMessages[114]) // TODO: change
}

func SetEmailEqualityError(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusBadRequest)
	r.SetError(114, errCodesMessages[114]) // TODO: change
}

func SetPhoneEqualityError(r *core_dtos.ResponseDto) {
	r.SetStatus(fiber.StatusBadRequest)
	r.SetError(113, errCodesMessages[113]) // TODO: change
}

func SetInvalidCredentialError(r *core_dtos.ResponseDto) {
	r.SetError(113, errCodesMessages[113]) // TODO: change
}
