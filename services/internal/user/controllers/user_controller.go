package controllers

import (
	"log"
	"services/internal/user/dto"
	"services/internal/user/services"
	"services/pkg/common/response"
	"services/pkg/common/validation"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service_user    *services.UserService
	service_session *services.SessionService
	service_account *services.AccountService
	service_socket  *services.SocketService
	validator       validation.Valitator
}

func NewUserController(user *services.UserService, session *services.SessionService, account *services.AccountService, socket services.SocketService) (*UserController, error) {
	validator := validation.NewValidator()
	return &UserController{
		service_user:    user,
		service_session: session,
		service_account: account,
		validator:       *validator,
		service_socket:  &socket,
	}, nil
}

// RegisterUserWithLocalParameters
// @Summary Register With Given Credentials.
// @Description Register Endpoint.
// @Produce json
// @Param body body dto.RegisterUserRequestDTO true "Enter Credential."
// @Success 201 {object} dto.SwaggerSuccessResponse[any,any] "ok" "返回用户信息"
// @Failure 400 {string} string "err_code：10002 参数错误； err_code：10003 校验错误"
// @Failure 401 {string} string "err_code：10001 登录失败"
// @Failure 500 {string} string "err_code：20001 服务错误；err_code：20002 接口错误；err_code：20003 无数据错误；err_code：20004 数据库异常；err_code：20005 缓存异常"
// @Router /register [post]
func (c *UserController) RegisterUserWithLocalParameters(ctx *fiber.Ctx) error {
	userReq, err := dto.NewRegisterUserRequestDTO(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Body Parse Error.", err.Error())
	}
	err = c.validator.Validate(userReq)
	if err != nil {
		return err
	}
	// Create User with Data
	user, err := c.service_user.CreateUserWithData(userReq.Name, userReq.LastName, userReq.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Model Creation Error.", err.Error())
	}

	/// Create New Account and Connect with User.
	_, err = c.service_account.ConnectWithUser(user, userReq.Username, userReq.Password)
	if err != nil {
		// If has an any exception, delete user that given user id.
		err2 := c.service_user.DeleteUserWithID(user.ID)
		if err2 != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Model Creation Error.There are To Many Problems.", err.Error())
		}
		//
		return fiber.NewError(fiber.StatusInternalServerError, "User Didn Connect With Account.", err.Error())

	}
	///
	return response.ResponseWithSuccessMessage(ctx, fiber.StatusCreated, "Account Creation Process Successful.", nil, nil)
}

// LoginUserWithLocalParameters
// @Summary Login With Given Credentials.
// @Description test deneme
// @Produce json
// @Param body body dto.LoginUserRequestDto true "body参数"
// @Success 200 {object} dto.LoginUserResponseDTO "ok" "返回用户信息"
// @Failure 400 {string} string "err_code：10002 参数错误； err_code：10003 校验错误"
// @Failure 401 {string} string "err_code：10001 登录失败"
// @Failure 500 {string} string "err_code：20001 服务错误；err_code：20002 接口错误；err_code：20003 无数据错误；err_code：20004 数据库异常；err_code：20005 缓存异常"
// @Router /login [post]
func (c *UserController) LoginUserWithLocalParameters(ctx *fiber.Ctx) error {
	credential, err := dto.NewLoginUserRequestDTO(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Given Credentials Wrong.", err.Error())
	}

	err = c.validator.Validate(credential)
	if err != nil {
		return err
	}
	account, err := c.service_account.GetAccountWithUsername(credential.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Given Credential is Wrong. Please Check And Try Again.", err.Error())
	}
	err = c.service_account.CompareHashAndPassword(account.Password, credential.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User Not Found or Wrong Password.")
	}
	// Pass Password Chech and User Exist.

	session, err := c.service_session.CreateSession(account)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "Session Creation Error.", err.Error())
	}
	dto := dto.NewLoginUserResponseDTO(session.AccessToken, session.RefreshToken)

	return response.ResponseWithSuccessMessage(ctx, fiber.StatusCreated, "Login SuccessFul.", dto, nil)
}

// WhoAmI
// @Summary Gives Who Am I.
// @Description @Me Endpoint.
// @Produce json
// @Success 200 {object} dto.SwaggerSuccessResponse[dto.WhoAmIResponseDto,any] "ok" "返回用户信息"
// @Failure 400 {string} string "err_code：10002 参数错误； err_code：10003 校验错误"
// @Failure 401 {string} string "err_code：10001 登录失败"
// @Failure 500 {string} string "err_code：20001 服务错误；err_code：20002 接口错误；err_code：20003 无数据错误；err_code：20004 数据库异常；err_code：20005 缓存异常"
// @Router /me [get]
// @Security BearerAuth
func (c *UserController) GetUserWithToken(ctx *fiber.Ctx) error {

	body, err := dto.NewWhoAmIRequestDTO(ctx)

	if err != nil {
		return err
	}

	account, err := c.service_session.GetAccountWithAccessToken(body.AccessToken)

	if err != nil {
		return err
	}

	resp := dto.NewWhoAmIResponseDTO(*account)

	return response.ResponseWithSuccessMessage(ctx, fiber.StatusAccepted, "WhoAmI", resp, nil)
}

// RefreshTokens
// @Summary Gives All Tokens with Refresh.
// @Description @Me Endpoint.
// @Produce json
// @Success 200 {object} dto.SwaggerSuccessResponse[dto.WhoAmIResponseDto,any] "ok" "返回用户信息"
// @Failure 400 {string} string "err_code：10002 参数错误； err_code：10003 校验错误"
// @Failure 401 {string} string "err_code：10001 登录失败"
// @Failure 500 {string} string "err_code：20001 服务错误；err_code：20002 接口错误；err_code：20003 无数据错误；err_code：20004 数据库异常；err_code：20005 缓存异常"
// @Router /me [get]
// @Security BearerAuth
func (c *UserController) GetTokensWithRefreshToken(ctx *fiber.Ctx) error {

	return nil
}

func (c *UserController) GetWebSocketHandler(conn *websocket.Conn) {
	// c.Locals is added to the *websocket.Conn
	log.Println(conn.Locals("allowed"))  // true
	log.Println(conn.Params("id"))       // 123
	log.Println(conn.Query("v"))         // 1.0
	log.Println(conn.Cookies("session")) // ""

	for {

		_, msg := c.service_socket.ReadMessage(conn)

		parsed := c.service_socket.ParseMessage(conn, msg)

		switch parsed.Type {
		case "ping":
			c.service_socket.WriteMessage(conn, websocket.TextMessage, []byte("Pong"))

		default:
			c.service_socket.WriteMessage(conn, websocket.TextMessage, []byte("Wrong Message. Please Check Again."))
		}
	}

}
