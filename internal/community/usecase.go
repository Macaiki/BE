package community

import "macaiki/internal/community/dto"

type CommunityUsecase interface {
	GetAllCommunities(userID int, search string) ([]dto.CommunityDetailResponse, error)
	GetCommunity(userID, communityID uint) (dto.CommunityDetailResponse, error)
	StoreCommunity(community dto.CommunityRequest, role string) error
	UpdateCommunity(id uint, community dto.CommunityRequest, role string) (dto.CommunityResponse, error)
	DeleteCommunity(id uint, role string) error

	FollowCommunity(userID, communityID uint) error
	UnfollowCommunity(userID, communityID uint) error
}
