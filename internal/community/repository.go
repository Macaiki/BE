package community

import (
	communityEntity "macaiki/internal/community/entity"
	userEntity "macaiki/internal/user/entity"
)

type CommunityRepository interface {
	GetAllCommunities(userID uint, search string) ([]communityEntity.Community, error)
	GetCommunityWithDetail(userID, communityID uint) (communityEntity.Community, error)
	GetCommunity(id uint) (communityEntity.Community, error)
	StoreCommunity(community communityEntity.Community) error
	UpdateCommunity(community communityEntity.Community, communityReq communityEntity.Community) error
	DeleteCommunity(community communityEntity.Community) error

	FollowCommunity(user userEntity.User, community communityEntity.Community) error
	UnfollowCommunity(user userEntity.User, community communityEntity.Community) error
}