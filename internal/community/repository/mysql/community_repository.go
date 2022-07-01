package mysql

import (
	"macaiki/internal/community"
	communityEntity "macaiki/internal/community/entity"
	userEntity "macaiki/internal/user/entity"

	"gorm.io/gorm"
)

type CommunityRepositoryImpl struct {
	db *gorm.DB
}

func NewCommunityRepository(db *gorm.DB) community.CommunityRepository {
	return &CommunityRepositoryImpl{db}
}

func (cr *CommunityRepositoryImpl) GetAllCommunity(search string) ([]communityEntity.Community, error) {
	communities := []communityEntity.Community{}

	res := cr.db.Where("name LIKE ?", "%"+search+"%").Find(&communities)
	err := res.Error
	if err != nil {
		return []communityEntity.Community{}, err
	}

	return communities, nil
}

func (cr *CommunityRepositoryImpl) GetAllCommunityDetail(userID, search string) ([]communityEntity.CommunityWithDetail, error) {
	communitiesWithDetail := []communityEntity.CommunityWithDetail{}

	res := cr.db.Raw("SELECT c.*, !ISNULL(cf.user_id) AS `is_followed` FROM `communities` AS c LEFT JOIN (SELECT * FROM community_followers WHERE user_id = ?) AS cf ON c.id = cf.community_id", userID).Scan(&communitiesWithDetail)
	err := res.Error
	if err != nil {
		return []communityEntity.CommunityWithDetail{}, err
	}

	return communitiesWithDetail, nil
}

func (cr *CommunityRepositoryImpl) GetCommunity(id uint) (communityEntity.Community, error) {
	community := communityEntity.Community{}

	res := cr.db.Find(&community, id)
	err := res.Error

	if err != nil {
		return communityEntity.Community{}, err
	}

	return community, nil
}

func (cr *CommunityRepositoryImpl) StoreCommunity(community communityEntity.Community) error {
	res := cr.db.Create(&community)
	err := res.Error
	if err != nil {
		return err
	}

	return nil
}

func (cr *CommunityRepositoryImpl) UpdateCommunity(community communityEntity.Community, communityReq communityEntity.Community) error {

	res := cr.db.Model(&community).Updates(communityReq)
	err := res.Error
	if err != nil {
		return err
	}

	return nil
}

func (cr *CommunityRepositoryImpl) DeleteCommunity(community communityEntity.Community) error {
	res := cr.db.Delete(&community)
	err := res.Error
	if err != nil {
		return err
	}
	return nil
}
func (cr *CommunityRepositoryImpl) FollowCommunity(user userEntity.User, community communityEntity.Community) error {
	err := cr.db.Model(&community).Association("Users").Append(&user)
	if err != nil {
		return err
	}

	return nil
}

func (cr *CommunityRepositoryImpl) UnfollowCommunity(user userEntity.User, community communityEntity.Community) error {
	err := cr.db.Model(&community).Association("Users").Delete(&user)
	if err != nil {
		return err
	}

	return nil
}
