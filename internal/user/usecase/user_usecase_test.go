package usecase

import (
	notifEntity "macaiki/internal/notification/entity"
	notifMock "macaiki/internal/notification/mocks"
	userDTO "macaiki/internal/user/dto"
	userEntity "macaiki/internal/user/entity"
	userMock "macaiki/internal/user/mocks"
	"macaiki/pkg/utils"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	mockUserEntity1 = userEntity.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:              "dummy@gmail.com",
		Username:           "dummy",
		Password:           "$2a$04$UP.ZNuepVAiEedwlZrvA3.ywqqNszceSuqnZQl4mozYOzO9ILY2kK",
		Name:               "dummy",
		ProfileImageUrl:    "dummy",
		BackgroundImageUrl: "dummy",
		Bio:                "dummy",
		Profession:         "dummy",
		Role:               "User",
		EmailVerifiedAt:    time.Now(),
		IsBanned:           0,
	}

	mockUserEntity2 = userEntity.User{
		Model: gorm.Model{
			ID:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:              "dummy@gmail.com",
		Username:           "dummy",
		Password:           "123456",
		Name:               "dummy",
		ProfileImageUrl:    "dummy",
		BackgroundImageUrl: "dummy",
		Bio:                "dummy",
		Profession:         "dummy",
		Role:               "Admin",
		EmailVerifiedAt:    time.Now(),
		IsBanned:           0,
	}

	mockedUserArr = []userEntity.User{mockUserEntity1, mockUserEntity1, mockUserEntity1}

	v = validator.New()
)

// func TestLogin(t *testing.T) {
// 	mockUserRepo := userMock.NewUserRepository(t)

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("GetByEmail", mockUserEntity1.Email).Return(userEntity.User{}, nil).Once()

// 		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)
// 	})
// }

// func TestRegister(t *testing.T) {
// 	mockUserRepo := userMock.NewUserRepository(t)

// 	mockUserReq := userDTO.UserRequest{
// 		Email:                "dummy@gmail.com",
// 		Username:             "dummy",
// 		Password:             "123456",
// 		PasswordConfirmation: "123456",
// 	}

// 	t.Run("success", func(t *testing.T) {
// 		mockUserRepo.On("GetByEmail", mockUserReq.Email).Return(userEntity.User{}, nil).Once()
// 		mockUserRepo.On("GetByUsername", mockUserReq.Username).Return(userEntity.User{}, nil).Once()
// 		mockUserRepo.On("Store", mockUserEntity1).Return(nil).Once()

// 		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)
// 		err := testUserUsecase.Register(mockUserReq)

// 		assert.NoError(t, err)
// 	})
// }

func TestGetAll(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetAllWithDetail", uint(1), "").Return(mockedUserArr, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.GetAll(uint(1), "")

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("GetAllWithDetail", uint(1), "").Return(mockedUserArr, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.GetAll(uint(1), "")

		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGet(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetWithDetail", uint(1), uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollowingNumber", uint(1)).Return(10, nil).Once()
		mockUserRepo.On("GetFollowerNumber", uint(1)).Return(10, nil).Once()
		mockUserRepo.On("GetThreadsNumber", uint(1)).Return(10, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.Get(uint(1), uint(1))

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mockUserRepo.On("GetWithDetail", uint(1), uint(1)).Return(userEntity.User{}, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.Get(uint(1), uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("internal-server-error-1", func(t *testing.T) {
		mockUserRepo.On("GetWithDetail", uint(1), uint(1)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.Get(uint(1), uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("internal-server-error-2", func(t *testing.T) {
		mockUserRepo.On("GetWithDetail", uint(1), uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollowingNumber", uint(1)).Return(0, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.Get(uint(1), uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("internal-server-error-3", func(t *testing.T) {
		mockUserRepo.On("GetWithDetail", uint(1), uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollowingNumber", uint(1)).Return(10, nil).Once()
		mockUserRepo.On("GetFollowerNumber", uint(1)).Return(10, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.Get(uint(1), uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("internal-server-error-4", func(t *testing.T) {
		mockUserRepo.On("GetWithDetail", uint(1), uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollowingNumber", uint(1)).Return(10, nil).Once()
		mockUserRepo.On("GetFollowerNumber", uint(1)).Return(10, nil).Once()
		mockUserRepo.On("GetThreadsNumber", uint(1)).Return(10, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)
		res, err := testUserUsecase.Get(uint(1), uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestUpdate(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)

	mockUserUpdateDTO := userDTO.UserUpdateRequest{
		Name:       "dummy",
		Bio:        "dummy",
		Profession: "dummy",
	}

	mockUserEntityUpdate := userEntity.User{
		Name:       "dummy",
		Bio:        "dummy",
		Profession: "dummy",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Update", &mockUserEntity1, mockUserEntityUpdate).Return(mockUserEntity1, nil).Once()

		testUseUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		res, err := testUseUsecase.Update(mockUserUpdateDTO, uint(1))

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, nil).Once()

		testUseUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		res, err := testUseUsecase.Update(mockUserUpdateDTO, uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUseUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		res, err := testUseUsecase.Update(mockUserUpdateDTO, uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Update", &mockUserEntity1, mockUserEntityUpdate).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUseUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		res, err := testUseUsecase.Update(mockUserUpdateDTO, uint(1))

		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestDelete(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Delete", uint(1)).Return(nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		err := testUserUsecase.Delete(uint(1), uint(1), "Admin")

		assert.NoError(t, err)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		err := testUserUsecase.Delete(uint(1), uint(1), "Admin")

		assert.Error(t, err)
	})

	t.Run("unautorize", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		err := testUserUsecase.Delete(uint(1), uint(2), "User")

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Delete", uint(1)).Return(utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		err := testUserUsecase.Delete(uint(1), uint(1), "Admin")

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		err := testUserUsecase.Delete(uint(1), uint(1), "Admin")

		assert.Error(t, err)
	})
}

func TestChangeEmail(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)

	mockInfoDTOReqSuccess := userDTO.UserLoginRequest{
		Email:    "dummyupdate@gmail.com",
		Password: "123456",
	}

	mockInfoDTOReqSuccessFail1 := userDTO.UserLoginRequest{
		Email:    "",
		Password: "123456",
	}

	mockInfoDTOReqSuccessFail2 := userDTO.UserLoginRequest{
		Email:    "dummyupdate@gmail.com",
		Password: "1234567",
	}

	mockEntityReq := userEntity.User{
		Email: "dummyupdate@gmail.com",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetByEmail", mockInfoDTOReqSuccess.Email).Return(userEntity.User{}, nil).Once()
		mockUserRepo.On("Update", &mockUserEntity1, mockEntityReq).Return(mockUserEntity1, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)

		res, err := testUserUsecase.ChangeEmail(uint(1), mockInfoDTOReqSuccess)

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)

		res, err := testUserUsecase.ChangeEmail(uint(1), mockInfoDTOReqSuccess)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetByEmail", mockInfoDTOReqSuccess.Email).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)

		res, err := testUserUsecase.ChangeEmail(uint(1), mockInfoDTOReqSuccess)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetByEmail", mockInfoDTOReqSuccess.Email).Return(userEntity.User{}, nil).Once()
		mockUserRepo.On("Update", &mockUserEntity1, mockEntityReq).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)

		res, err := testUserUsecase.ChangeEmail(uint(1), mockInfoDTOReqSuccess)

		assert.Error(t, err)
		assert.Empty(t, res)
	})
	t.Run("bad-param-input", func(t *testing.T) {
		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)

		res, err := testUserUsecase.ChangeEmail(uint(1), mockInfoDTOReqSuccessFail1)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("email-already-used", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetByEmail", mockInfoDTOReqSuccess.Email).Return(mockUserEntity1, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)

		res, err := testUserUsecase.ChangeEmail(uint(1), mockInfoDTOReqSuccess)

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetByEmail", mockInfoDTOReqSuccess.Email).Return(userEntity.User{}, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, v, nil, nil)

		res, err := testUserUsecase.ChangeEmail(uint(1), mockInfoDTOReqSuccessFail2)

		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestChangePassword(t *testing.T) {}

func TestGetUserFollowers(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollower", uint(1), uint(1)).Return(mockedUserArr, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowers(uint(1), uint(1))

		assert.NoError(t, err)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowers(uint(1), uint(1))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowers(uint(1), uint(1))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollower", uint(1), uint(1)).Return([]userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowers(uint(1), uint(1))

		assert.Error(t, err)
	})
}

func TestGetUserFollowing(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollowing", uint(1), uint(1)).Return(mockedUserArr, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowing(uint(1), uint(1))

		assert.NoError(t, err)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowing(uint(1), uint(1))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowing(uint(1), uint(1))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("GetFollowing", uint(1), uint(1)).Return([]userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, nil, nil, nil, nil, nil)

		_, err := testUserUsecase.GetUserFollowing(uint(1), uint(1))

		assert.Error(t, err)
	})
}

// func TestSetProfileImage(t *testing.T) {}

// func TestSetBackgroundImage(t *testing.T) {}

func TestFollow(t *testing.T) {
	mockUserRepo := userMock.NewUserRepository(t)
	mockNotifRepo := notifMock.NewNotificationRepository(t)

	mockNotifEntity := notifEntity.Notification{
		UserID:            1,
		NotificationType:  "Follow You",
		NotificationRefID: 2,
		IsReaded:          0,
	}
	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Get", uint(2)).Return(mockUserEntity2, nil).Once()
		mockUserRepo.On("Follow", mockUserEntity1, mockUserEntity2).Return(mockUserEntity1, nil).Once()
		mockNotifRepo.On("StoreNotification", mockNotifEntity).Return(nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(2))

		assert.NoError(t, err)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(2))

		assert.Error(t, err)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Get", uint(2)).Return(userEntity.User{}, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(2))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(2))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Get", uint(2)).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(2))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Get", uint(2)).Return(mockUserEntity2, nil).Once()
		mockUserRepo.On("Follow", mockUserEntity1, mockUserEntity2).Return(userEntity.User{}, utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(2))

		assert.Error(t, err)
	})

	t.Run("internal-server-error", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Get", uint(2)).Return(mockUserEntity2, nil).Once()
		mockUserRepo.On("Follow", mockUserEntity1, mockUserEntity2).Return(mockUserEntity1, nil).Once()
		mockNotifRepo.On("StoreNotification", mockNotifEntity).Return(utils.ErrInternalServerError).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(2))

		assert.NoError(t, err)
	})

	t.Run("bad-param-input", func(t *testing.T) {
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()
		mockUserRepo.On("Get", uint(1)).Return(mockUserEntity1, nil).Once()

		testUserUsecase := NewUserUsecase(mockUserRepo, nil, nil, mockNotifRepo, nil, nil, nil, nil)

		err := testUserUsecase.Follow(uint(1), uint(1))

		assert.Error(t, err)
	})
}
