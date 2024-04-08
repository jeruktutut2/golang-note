package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"golang-note/exception"
	"golang-note/globals"
	"golang-note/helper"
	modelentity "golang-note/model/entity"
	modelrequest "golang-note/model/request"
	modelresponse "golang-note/model/response"
	"golang-note/repository"
	"golang-note/util"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, requestId string, loginRequest modelrequest.LoginRequest) (loginResponse modelresponse.LoginResponse, accessToken string, refrshToken string, err error)
	RefreshToken(ctx context.Context, requestId string, username string, refreshToken string) (accessToken string, err error)
	LoginRedis(ctx context.Context, requestId string, loginRequest modelrequest.LoginRequest) (loginResponse modelresponse.LoginResponse, token string, err error)
	LoginMap(ctx context.Context, requestId string, loginRequest modelrequest.LoginRequest) (loginResponse modelresponse.LoginResponse, token string, err error)
}

type UserServiceImplementation struct {
	MysqlUtil                 util.MysqlUtil
	RedisUtil                 util.RedisUtil
	Validate                  *validator.Validate
	JwtKey                    string
	JwtAccessTokenExpireTime  uint16
	JwtRefreshTokenExpireTime uint16
	UserRepository            repository.UserRepository
	PermissionRepository      repository.PermissionRepository
	UserPermissionRepository  repository.UserPermissionRepository
	RedisRepository           repository.RedisRepository
}

func NewUserService(mysqlUtil util.MysqlUtil, redisUtil util.RedisUtil, validate *validator.Validate, jwtKey string, jwtAccessTokenExpireTime uint16, jwtRefreshTokenExpireTime uint16, userRepository repository.UserRepository, permissionRepository repository.PermissionRepository, userPermissionRepository repository.UserPermissionRepository, redisRepository repository.RedisRepository) UserService {
	return &UserServiceImplementation{
		MysqlUtil:                 mysqlUtil,
		RedisUtil:                 redisUtil,
		Validate:                  validate,
		JwtKey:                    jwtKey,
		JwtAccessTokenExpireTime:  jwtAccessTokenExpireTime,
		JwtRefreshTokenExpireTime: jwtRefreshTokenExpireTime,
		UserRepository:            userRepository,
		PermissionRepository:      permissionRepository,
		UserPermissionRepository:  userPermissionRepository,
		RedisRepository:           redisRepository,
	}
}

func (service *UserServiceImplementation) Login(ctx context.Context, requestId string, loginRequest modelrequest.LoginRequest) (loginResponse modelresponse.LoginResponse, accessToken string, refreshToken string, err error) {
	err = service.Validate.Struct(loginRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, loginRequest)
		if validationResult != nil {
			var validationResultByte []byte
			validationResultByte, err = json.Marshal(validationResult)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				err = exception.CheckError(err)
				return
			}
			err = exception.NewValidationException(string(validationResultByte))
			return
		}
	}

	tx, err := service.MysqlUtil.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	defer func() {
		errCommitOrRollback := service.MysqlUtil.CommitOrRollback(tx, err)
		if errCommitOrRollback != nil {
			helper.PrintLogToTerminal(err, requestId)
			loginResponse = modelresponse.LoginResponse{}
			err = exception.NewInternalServerErrorException()
		}
	}()

	var user *modelentity.User = &modelentity.User{}
	user.Username = sql.NullString{Valid: true, String: loginRequest.Username}
	err = service.UserRepository.GetByUsernameForUpdate(tx, ctx, user)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		// fmt.Println("err:", err)
		return
	} else if err == sql.ErrNoRows {
		err = exception.NewNotFoundException("cannot find user: " + loginRequest.Username)
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(loginRequest.Password))
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		if err != bcrypt.ErrMismatchedHashAndPassword {
			err = exception.CheckError(err)
			return
		} else {
			err = exception.NewBadRequestException("wrong username or password")
			return
		}
	}

	var userPermissions *[]modelentity.UserPermission = &[]modelentity.UserPermission{}
	err = service.UserPermissionRepository.GetByUserId(tx, ctx, user.Id.Int32, userPermissions)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	var permissionAccessToken []string
	if len(*userPermissions) > 0 {
		var ids []interface{}
		for _, userPermission := range *userPermissions {
			ids = append(ids, userPermission.Id.Int32)
		}

		var permissions *[]modelentity.Permission = &[]modelentity.Permission{}
		err = service.PermissionRepository.GetByInId(tx, ctx, ids, permissions)
		if err != nil {
			helper.PrintLogToTerminal(err, requestId)
			err = exception.CheckError(err)
			return
		}

		for _, permission := range *permissions {
			isContain := slices.Contains(permissionAccessToken, permission.Permission.String)
			// isContain := helper.CheckArrayContain(permissionAccessToken, permission.Permission.String)
			if !isContain {
				permissionAccessToken = append(permissionAccessToken, permission.Permission.String)
			}
		}
	}

	accessToken, err = generateAccessToken(*user, permissionAccessToken, service.JwtKey, service.JwtAccessTokenExpireTime)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		accessToken = ""
		refreshToken = ""
		err = exception.CheckError(err)
		return
	}
	refreshToken, err = generateRefreshToken(service.JwtKey, service.JwtRefreshTokenExpireTime)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		accessToken = ""
		refreshToken = ""
		err = exception.CheckError(err)
		return
	}

	rowsAffected, err := service.UserRepository.UpdateRefreshToken(tx, ctx, user.Id.Int32, refreshToken)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		accessToken = ""
		refreshToken = ""
		err = exception.CheckError(err)
		return
	}
	if rowsAffected != 1 {
		err = errors.New("rows affected not one: " + strconv.Itoa(int(rowsAffected)))
		accessToken = ""
		refreshToken = ""
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return
	}
	loginResponse = modelresponse.ToLoginResponse(*user)
	return
}

func generateAccessToken(user modelentity.User, permissionAccessToken []string, jwtKey string, jwtAccessTokenExpireTime uint16) (string, error) {
	key := []byte(jwtKey)
	now := time.Now()

	type AccessTokenClaims struct {
		Id          int32    `json:"id"`
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Permissions []string `json:"permissions"`
		jwt.RegisteredClaims
	}

	// Create the Claims
	accessTokenclaims := AccessTokenClaims{
		user.Id.Int32,
		user.Username.String,
		user.Email.String,
		permissionAccessToken,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(jwtAccessTokenExpireTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "login",
			// Subject:   "somebody",
			// ID:        "1",
			// Audience:  []string{"somebody_else"},
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenclaims)
	return jwtToken.SignedString(key)
}

func generateRefreshToken(jwtKey string, jwtRefreshTokenexpireTime uint16) (string, error) {
	key := []byte(jwtKey)
	now := time.Now()

	type RefreshTokenClaims struct {
		jwt.RegisteredClaims
	}

	// Create the Claims
	refreshTokenclaims := RefreshTokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.AddDate(0, 0, int(jwtRefreshTokenexpireTime))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "login",
			// Subject:   "somebody",
			// ID:        "1",
			// Audience:  []string{"somebody_else"},
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenclaims)
	return jwtToken.SignedString(key)
}

func (service *UserServiceImplementation) RefreshToken(ctx context.Context, requestId string, username string, refreshToken string) (accessToken string, err error) {
	tx, err := service.MysqlUtil.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	defer func() {
		errCommitOrRollback := service.MysqlUtil.CommitOrRollback(tx, err)
		if errCommitOrRollback != nil {
			helper.PrintLogToTerminal(err, requestId)
			accessToken = ""
			err = exception.NewInternalServerErrorException()
		}
	}()

	var user *modelentity.User = &modelentity.User{}
	user.RefreshToken = sql.NullString{Valid: true, String: refreshToken}
	err = service.UserRepository.GetByRefreshToken(tx, ctx, user)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	} else if err == sql.ErrNoRows {
		err = exception.NewNotFoundException("cannot find user by refresh token")
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	var userPermissions *[]modelentity.UserPermission = &[]modelentity.UserPermission{}
	err = service.UserPermissionRepository.GetByUserId(tx, ctx, user.Id.Int32, userPermissions)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	var permissionAccessToken []string
	if len(*userPermissions) > 0 {
		var ids []interface{}
		for _, userPermission := range *userPermissions {
			ids = append(ids, userPermission.Id.Int32)
		}

		var permissions *[]modelentity.Permission = &[]modelentity.Permission{}
		err = service.PermissionRepository.GetByInId(tx, ctx, ids, permissions)
		if err != nil {
			helper.PrintLogToTerminal(err, requestId)
			err = exception.CheckError(err)
			return
		}

		for _, permission := range *permissions {
			isContain := slices.Contains(permissionAccessToken, permission.Permission.String)
			if !isContain {
				permissionAccessToken = append(permissionAccessToken, permission.Permission.String)
			}
		}
	}

	accessToken, err = generateAccessToken(*user, permissionAccessToken, service.JwtKey, service.JwtAccessTokenExpireTime)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		accessToken = ""
		err = exception.CheckError(err)
		return
	}
	return
}

func (service *UserServiceImplementation) LoginRedis(ctx context.Context, requestId string, loginRequest modelrequest.LoginRequest) (loginResponse modelresponse.LoginResponse, token string, err error) {
	err = service.Validate.Struct(loginRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, loginRequest)
		if validationResult != nil {
			var validationResultByte []byte
			validationResultByte, err = json.Marshal(validationResult)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				err = exception.CheckError(err)
				return
			}
			err = exception.NewValidationException(string(validationResultByte))
			return
		}
	}

	tx, err := service.MysqlUtil.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	defer func() {
		errCommitOrRollback := service.MysqlUtil.CommitOrRollback(tx, err)
		if errCommitOrRollback != nil {
			helper.PrintLogToTerminal(err, requestId)
			loginResponse = modelresponse.LoginResponse{}
			err = exception.NewInternalServerErrorException()
		}
	}()

	var user *modelentity.User = &modelentity.User{}
	user.Username = sql.NullString{Valid: true, String: loginRequest.Username}
	err = service.UserRepository.GetByUsernameForUpdate(tx, ctx, user)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}
	if err == sql.ErrNoRows {
		err = exception.NewNotFoundException("cannot find user: " + loginRequest.Username)
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(loginRequest.Password))
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		if err != bcrypt.ErrMismatchedHashAndPassword {
			err = exception.CheckError(err)
			return
		} else {
			err = exception.NewBadRequestException("wrong username or password")
			return
		}
	}

	var userPermissions *[]modelentity.UserPermission = &[]modelentity.UserPermission{}
	err = service.UserPermissionRepository.GetByUserId(tx, ctx, user.Id.Int32, userPermissions)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	var permissionToken []string
	if len(*userPermissions) > 0 {
		var ids []interface{}
		for _, userPermission := range *userPermissions {
			ids = append(ids, userPermission.Id.Int32)
		}

		var permissions *[]modelentity.Permission = &[]modelentity.Permission{}
		err = service.PermissionRepository.GetByInId(tx, ctx, ids, permissions)
		if err != nil {
			helper.PrintLogToTerminal(err, requestId)
			err = exception.CheckError(err)
			return
		}

		for _, permission := range *permissions {
			isContain := slices.Contains(permissionToken, permission.Permission.String)
			if !isContain {
				permissionToken = append(permissionToken, permission.Permission.String)
			}
		}
	}

	token = uuid.New().String()
	permissionTokenString := strings.Join(permissionToken, ",")
	_, err = service.RedisRepository.Set(service.RedisUtil.GetClient(), ctx, token, permissionTokenString, 0)
	if err != nil && err != redis.Nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	loginResponse = modelresponse.ToLoginResponse(*user)
	return
}

func (service *UserServiceImplementation) LoginMap(ctx context.Context, requestId string, loginRequest modelrequest.LoginRequest) (loginResponse modelresponse.LoginResponse, token string, err error) {
	err = service.Validate.Struct(loginRequest)
	if err != nil {
		validationResult := helper.GetValidatorError(err, loginRequest)
		if validationResult != nil {
			var validationResultByte []byte
			validationResultByte, err = json.Marshal(validationResult)
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				err = exception.CheckError(err)
				return
			}
			err = exception.NewValidationException(string(validationResultByte))
			return
		}
	}

	tx, err := service.MysqlUtil.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	defer func() {
		errCommitOrRollback := service.MysqlUtil.CommitOrRollback(tx, err)
		if errCommitOrRollback != nil {
			helper.PrintLogToTerminal(err, requestId)
			loginResponse = modelresponse.LoginResponse{}
			err = exception.NewInternalServerErrorException()
		}
	}()

	var user *modelentity.User = &modelentity.User{}
	user.Username = sql.NullString{Valid: true, String: loginRequest.Username}
	err = service.UserRepository.GetByUsername(tx, ctx, user)
	if err != nil && err != sql.ErrNoRows {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	} else if err == sql.ErrNoRows {
		err = exception.NewNotFoundException("cannot find user: " + loginRequest.Username)
		helper.PrintLogToTerminal(err, requestId)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(loginRequest.Password))
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		if err != bcrypt.ErrMismatchedHashAndPassword {
			err = exception.CheckError(err)
			return
		} else {
			err = exception.NewBadRequestException("wrong username or password")
			return
		}
	}

	var userPermissions *[]modelentity.UserPermission = &[]modelentity.UserPermission{}
	err = service.UserPermissionRepository.GetByUserId(tx, ctx, user.Id.Int32, userPermissions)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.CheckError(err)
		return
	}

	var permissionToken []string
	if len(*userPermissions) > 0 {
		var ids []interface{}
		for _, userPermission := range *userPermissions {
			ids = append(ids, userPermission.Id.Int32)
		}

		var permissions *[]modelentity.Permission = &[]modelentity.Permission{}
		err = service.PermissionRepository.GetByInId(tx, ctx, ids, permissions)
		if err != nil {
			helper.PrintLogToTerminal(err, requestId)
			err = exception.CheckError(err)
			return
		}

		for _, permission := range *permissions {
			isContain := slices.Contains(permissionToken, permission.Permission.String)
			if !isContain {
				permissionToken = append(permissionToken, permission.Permission.String)
			}
		}
	}

	token = uuid.New().String()
	permissionTokenString := strings.Join(permissionToken, ",")
	globals.Session[token] = permissionTokenString

	loginResponse = modelresponse.ToLoginResponse(*user)
	return
}
