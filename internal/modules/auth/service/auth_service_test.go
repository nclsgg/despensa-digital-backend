package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/nclsgg/despensa-digital/backend/config"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/domain"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/model"
	"github.com/nclsgg/despensa-digital/backend/internal/modules/auth/service"
)

type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockAuthRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if usr, ok := args.Get(0).(*model.User); ok {
		return usr, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockAuthRepository) GetUserById(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, userID)
	if usr, ok := args.Get(0).(*model.User); ok {
		return usr, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockAuthRepository) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func getTestRedisClient(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return client, mr
}

func getTestConfig() *config.Config {
	return &config.Config{
		JWTExpiration: "1h",
		JWTIssuer:     "testIssuer",
		JWTAudience:   "testAudience",
		JWTSecret:     "secret",
	}
}

// Testes comentados pois os métodos foram removidos da interface AuthService
// após refatoração para usar apenas OAuth

func TestHashPassword(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	password := "minhasenha123"
	hashedPassword, err := authSvc.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
}

func TestAuthService_GetUserByEmail(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	user := &model.User{ID: uuid.New(), Email: "test@example.com"}
	repo.On("GetUser", mock.Anything, user.Email).Return(user, nil).Once()

	result, err := authSvc.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.Equal(t, user, result)

	repo.AssertExpectations(t)
}

func TestAuthService_GetUserByEmail_NotFound(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	email := "missing@example.com"
	repo.On("GetUser", mock.Anything, email).Return((*model.User)(nil), gorm.ErrRecordNotFound).Once()

	result, err := authSvc.GetUserByEmail(context.Background(), email)
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	require.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestAuthService_CreateUserOAuth_SetsPassword(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	user := &model.User{Email: "oauth@example.com"}
	repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once().Run(func(args mock.Arguments) {
		created := args.Get(1).(*model.User)
		require.NotEmpty(t, created.Password)
	})

	require.NoError(t, authSvc.CreateUserOAuth(context.Background(), user))
	repo.AssertExpectations(t)
}

func TestAuthService_CreateUserOAuth_Duplicate(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	user := &model.User{Email: "duplicate@example.com"}
	repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(gorm.ErrDuplicatedKey).Once()

	err := authSvc.CreateUserOAuth(context.Background(), user)
	require.ErrorIs(t, err, domain.ErrEmailAlreadyRegistered)

	repo.AssertExpectations(t)
}

func TestAuthService_CompleteProfile_Success(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	userID := uuid.New()
	existing := &model.User{ID: userID}
	repo.On("GetUserById", mock.Anything, userID).Return(existing, nil).Once()
	repo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once().Run(func(args mock.Arguments) {
		updated := args.Get(1).(*model.User)
		require.Equal(t, "John", updated.FirstName)
		require.Equal(t, "Doe", updated.LastName)
		require.True(t, updated.ProfileCompleted)
	})

	require.NoError(t, authSvc.CompleteProfile(context.Background(), userID, "John", "Doe"))

	repo.AssertExpectations(t)
}

func TestAuthService_CompleteProfile_NotFound(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	userID := uuid.New()
	repo.On("GetUserById", mock.Anything, userID).Return((*model.User)(nil), gorm.ErrRecordNotFound).Once()

	err := authSvc.CompleteProfile(context.Background(), userID, "John", "Doe")
	require.ErrorIs(t, err, domain.ErrUserNotFound)

	repo.AssertExpectations(t)
}

func TestAuthService_CompleteProfile_UpdateError(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	userID := uuid.New()
	existing := &model.User{ID: userID}
	repo.On("GetUserById", mock.Anything, userID).Return(existing, nil).Once()
	repo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(errors.New("db error")).Once()

	err := authSvc.CompleteProfile(context.Background(), userID, "John", "Doe")
	require.ErrorIs(t, err, domain.ErrProfileUpdateFailed)

	repo.AssertExpectations(t)
}

/*
func TestRegister_Success(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	user := &model.User{
		Email:    "teste@exemplo.com",
		Password: "senha123",
	}

	repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	accessToken, refreshToken, err := authSvc.Register(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEqual(t, "senha123", user.Password)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.True(t, mr.Exists("refresh_token:"+refreshToken))

	repo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	hashedPwd, err := authSvc.HashPassword("senhaCorreta")
	assert.NoError(t, err)

	user := &model.User{
		ID:       uuid.New(),
		Email:    "teste@exemplo.com",
		Password: hashedPwd,
	}

	repo.On("GetUser", mock.Anything, user.Email).Return(user, nil)

	accessToken, refreshToken, err := authSvc.Login(context.Background(), user.Email, "senhaErrada")
	assert.Error(t, err)
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	repo.AssertExpectations(t)
}

func TestLogout(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	userID := uuid.New()
	refreshToken, err := authSvc.GenerateRefreshToken(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	err = authSvc.Logout(context.Background(), refreshToken)
	assert.NoError(t, err)
	assert.False(t, mr.Exists("refresh_token:"+refreshToken))
}

func TestGenerateAccessToken(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	user := &model.User{
		ID:    uuid.New(),
		Email: "teste@exemplo.com",
	}

	tokenString, err := authSvc.GenerateAccessToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	parsedToken, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestRefreshToken(t *testing.T) {
	redisClient, mr := getTestRedisClient(t)
	defer mr.Close()

	repo := new(mockAuthRepository)
	cfg := getTestConfig()
	authSvc := service.NewAuthService(repo, cfg, redisClient)

	userID := uuid.New()
	user := &model.User{
		ID:    userID,
		Email: "teste@exemplo.com",
	}

	refreshToken, err := authSvc.GenerateRefreshToken(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	repo.On("GetUserById", mock.Anything, userID).Return(user, nil)

	newAccessToken, newRefreshToken, err := authSvc.RefreshToken(context.Background(), refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)

	repo.AssertExpectations(t)
}
*/
