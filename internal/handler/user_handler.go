package handler

import (
	"users-service/internal/dto"
	"users-service/internal/logger"
	"users-service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UserController struct {
	UserService *service.UserService
}

// @Summary Create new user
// @Description that method creates a new user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateDTO
// @Param user body dto.CreateUserDTO true "user data"
// @Failure 409 {object} dto.ErrorDTO
// @Failure 500 {object} dto.ErrorDTO
// @Failure 400 {object} dto.ErrorDTO
// @Router /users/create [post]
func (u *UserController) CreateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.CreateUserDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in create user user controller", zap.Error(err), zap.String("function", "usercontroller.createuser"))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		status, message := u.UserService.CreateUser(input, ctx.Context())
		logger.ZapLogger.Info("returning reply create user")
		if status >= 400 {
			return ctx.Status(status).JSON(fiber.Map{"error": message})
		}
		return ctx.Status(status).JSON(message)
	}
}

// @Summary Find One User
// @Description that method finds an user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} dto.FindUserDTO
// @Failure 404 {object} dto.ErrorDTO
// @Failure 500 {object} dto.ErrorDTO
// @Router /users/one/{id} [get]
func (u *UserController) FindOneUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		status, reply := u.UserService.FindOneUserById(id, ctx.Context())
		if status >= 400 {
			logger.ZapLogger.Error("error in find one user user controller", zap.Any("error", reply))
			return ctx.Status(status).JSON(fiber.Map{"error": reply})
		}
		logger.ZapLogger.Info("user was searched. returning reply")
		return ctx.Status(status).JSON(reply)
	}
}

// @Summary Login User
// @Description that is the login method
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.LoginUserDTO true "user data"
// @Success 200 {object} dto.LoginDTO
// @Failure 401 {object} dto.ErrorDTO
// @Failure 409 {object} dto.ErrorDTO
// @Failure 500 {object} dto.ErrorDTO
// @Router /users/login [post]
func (u *UserController) LoginUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.LoginUserDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in body parser", zap.Error(err), zap.String("function", "user controller login user"))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		status, reply := u.UserService.LoginUser(input, ctx.Context())
		if status >= 400 {
			return ctx.Status(status).JSON(fiber.Map{"error": reply})
		}
		return ctx.Status(status).JSON(reply)
	}
}

// @Summary Update User
// @Description update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Param user body dto.UpdateUserDTO true "user data"
// @Success 200 {object} dto.MessageDTO
// @Failure 400 {object} dto.ErrorDTO
// @Failure 404 {object} dto.ErrorDTO
// @Failure 500 {object} dto.ErrorDTO
// @Router /users/update/{id} [put]
func (u *UserController) UpdateOneUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		id := user["id"].(string)
		var input dto.UpdateUserDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in body parser", zap.Error(err), zap.String("function", "user controller.updateoneuser"))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		status, reply := u.UserService.UpdateUser(input, ctx.Context(), id)
		if status >= 400 {
			logger.ZapLogger.Error("error in user service update user", zap.String("function", "user controller.updateoneuser"))
			return ctx.Status(status).JSON(fiber.Map{"error": reply})
		}
		return ctx.Status(status).JSON(fiber.Map{"message": reply})
	}
}

// @Summary Update User Role
// @Description update user role
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Param user body dto.UpdateUserRoleDTO true "user data"
// @Success 200 {object} dto.MessageDTO
// @Failure 400 {object} dto.ErrorDTO
// @Failure 404 {object} dto.ErrorDTO
// @Failure 500 {object} dto.ErrorDTO
// @Router /users/update/role/{id} [patch]
func (u *UserController) UpdateOneUserRole() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		var input dto.UpdateUserRoleDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in body parser", zap.Error(err), zap.String("function", "usercontroller.updateoneuserrole"))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		status, reply := u.UserService.UpdateUserRole(input, ctx.Context(), id)
		if status >= 400 {
			logger.ZapLogger.Error("error in user service update user", zap.String("function", "user controller.updateoneuser"))
			return ctx.Status(status).JSON(fiber.Map{"error": reply})
		}
		return ctx.Status(status).JSON(fiber.Map{"message": reply})
	}
}


// @Summary Find All Users
// @Description find all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageDTO
// @Failure 500 {object} dto.ErrorDTO
// @Router /users/all [get]
func (u *UserController) FindAllUsers() fiber.Handler {
	return  func(ctx *fiber.Ctx) error {
		status, users := u.UserService.FindAllUsers()
		if status >= 400 {
			logger.ZapLogger.Error("error in user service find all users", zap.String("function", "user controller.findallusers"))
			return ctx.Status(status).JSON(fiber.Map{"error": "internal server error"})
		}
		return  ctx.Status(status).JSON(users)
	}
}

func (u *UserController) VerifyUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		var input dto.VerifyCodeDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in body parser", zap.Error(err))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		id := user["id"].(string)

		status, message := u.UserService.VerifyCode(id, ctx.Context(), input)
		if status >= 400 {
			return  ctx.Status(status).JSON(fiber.Map{"error": message})
		}
		return  ctx.Status(status).JSON(fiber.Map{"message": message})
	}
}