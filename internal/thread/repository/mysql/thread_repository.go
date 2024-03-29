package mysql

import (
	"fmt"
	"macaiki/internal/thread"
	"macaiki/internal/thread/entity"
	"macaiki/pkg/utils"
	"strings"

	"gorm.io/gorm"
)

type ThreadRepositoryImpl struct {
	db *gorm.DB
}

func CreateNewThreadRepository(db *gorm.DB) thread.ThreadRepository {
	return &ThreadRepositoryImpl{db: db}
}

func (tr *ThreadRepositoryImpl) SetThreadImage(imageURL string, threadID uint) error {
	fmt.Println(imageURL)
	res := tr.db.Model(&entity.Thread{}).Where("id = ?", threadID).Update("image_url", imageURL)

	if res.Error != nil {
		fmt.Println(res.Error)
		if res.Error.Error() == "record not found" {
			return utils.ErrNotFound
		}
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetThreadByID(threadID uint) (entity.Thread, error) {
	var thread entity.Thread
	res := tr.db.First(&thread, threadID)
	if res.Error != nil {
		fmt.Println(res.Error)
		if res.Error.Error() == "record not found" {
			return entity.Thread{}, utils.ErrNotFound
		}
		return entity.Thread{}, utils.ErrInternalServerError
	}

	return thread, nil
}

func (tr *ThreadRepositoryImpl) CreateThread(thread entity.Thread) (entity.Thread, error) {
	res := tr.db.Debug().Create(&thread)
	if res.Error != nil {
		fmt.Println(res.Error)
		return entity.Thread{}, utils.ErrInternalServerError
	}

	return thread, nil
}

func (tr *ThreadRepositoryImpl) DeleteThread(threadID uint) error {
	res := tr.db.Delete(&entity.Thread{}, threadID)
	if res.Error != nil {
		return utils.ErrInternalServerError
	}

	if res.RowsAffected < 1 {
		return utils.ErrNotFound
	}

	return nil
}

func (tr *ThreadRepositoryImpl) UpdateThread(threadID uint, thread entity.Thread) error {
	res := tr.db.Model(&entity.Thread{}).Where("id", threadID).Updates(thread)
	if res.Error != nil {
		fmt.Println(res.Error)
		if res.Error.Error() == "record not found" {
			return utils.ErrNotFound
		}
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) UpvoteThread(threadUpvote entity.ThreadUpvote) error {
	res := tr.db.Create(&threadUpvote)

	if res.Error != nil {
		fmt.Println(res.Error)
		if strings.HasPrefix(res.Error.Error(), "Error 1452: Cannot add or update a child row") {
			return utils.ErrNotFound
		} else if strings.HasPrefix(res.Error.Error(), "Error 1062: Duplicate entry") {
			return utils.ErrDuplicateEntry
		}
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetTrendingThreads(userID uint) ([]entity.ThreadWithDetails, error) {
	var threads []entity.ThreadWithDetails
	// TODO: retrieve name, profile URL, etc
	res := tr.db.Raw("SELECT t.*, t2.upvotes_count, NOT ISNULL(t3.id) AS is_upvoted, NOT ISNULL(t4.user_id) AS is_followed, NOT ISNULL(t5.id) AS is_downvoted, users.name, users.profile_image_url, users.profession FROM threads t LEFT JOIN (SELECT thread_id, COUNT(*) AS upvotes_count FROM thread_upvotes tu WHERE DATEDIFF(NOW(), tu.created_at) < 7 AND tu.deleted_at IS NULL GROUP BY thread_id) AS t2 ON t.id = t2.thread_id LEFT JOIN (SELECT * FROM thread_upvotes tu WHERE tu.user_id = ? AND tu.deleted_at IS NULL) AS t3 ON t.id = t3.thread_id LEFT JOIN (SELECT user_id FROM user_followers uf WHERE uf.follower_id = ?) AS t4 ON t4.user_id = t.user_id LEFT JOIN users ON users.id = t.user_id LEFT JOIN (SELECT id, thread_id FROM thread_downvotes td WHERE td.user_id = ? AND td.deleted_at IS NULL) AS t5 ON t5.thread_id = t.id WHERE t.deleted_at IS NULL ORDER BY upvotes_count DESC;", userID, userID, userID).Scan(&threads)
	if res.Error != nil {
		return []entity.ThreadWithDetails{}, res.Error
	}

	return threads, nil
}

func (tr *ThreadRepositoryImpl) GetTrendingThreadsWithLimit(userID uint, limit int) ([]entity.ThreadWithDetails, error) {
	var threads []entity.ThreadWithDetails
	// TODO: retrieve name, profile URL, etc
	res := tr.db.Raw("SELECT t.*, t2.upvotes_count, NOT ISNULL(t3.id) AS is_upvoted, NOT ISNULL(t4.user_id) AS is_followed, NOT ISNULL(t5.id) AS is_downvoted, users.name, users.profile_image_url, users.profession FROM threads t LEFT JOIN (SELECT thread_id, COUNT(*) AS upvotes_count FROM thread_upvotes tu WHERE DATEDIFF(NOW(), tu.created_at) < 7 AND tu.deleted_at IS NULL GROUP BY thread_id) AS t2 ON t.id = t2.thread_id LEFT JOIN (SELECT * FROM thread_upvotes tu WHERE tu.user_id = ? AND tu.deleted_at IS NULL) AS t3 ON t.id = t3.thread_id LEFT JOIN (SELECT user_id FROM user_followers uf WHERE uf.follower_id = ?) AS t4 ON t4.user_id = t.user_id LEFT JOIN users ON users.id = t.user_id LEFT JOIN (SELECT id, thread_id FROM thread_downvotes td WHERE td.user_id = ? AND td.deleted_at IS NULL) AS t5 ON t5.thread_id = t.id WHERE t.deleted_at IS NULL ORDER BY upvotes_count DESC LIMIT ?;", userID, userID, userID, limit).Scan(&threads)
	if res.Error != nil {
		return []entity.ThreadWithDetails{}, res.Error
	}

	return threads, nil
}

func (tr *ThreadRepositoryImpl) GetThreadsFromFollowedCommunity(userID uint) ([]entity.ThreadWithDetails, error) {
	var threads []entity.ThreadWithDetails

	res := tr.db.Raw("SELECT t.*, t2.upvotes_count, NOT ISNULL(t4.id) AS is_upvoted, NOT ISNULL(t5.user_id) AS is_followed, NOT ISNULL(t6.id) AS is_downvoted, users.name, users.profile_image_url, users.profession FROM threads t LEFT JOIN (SELECT thread_id, COUNT(*) AS upvotes_count FROM thread_upvotes tu WHERE tu.deleted_at IS NULL GROUP BY thread_id) AS t2 ON t.id = t2.thread_id INNER JOIN (SELECT * FROM community_followers cf WHERE cf.user_id = ?) AS t3 ON t.community_id = t3.community_id INNER JOIN users ON users.id = t.user_id LEFT JOIN (SELECT * FROM thread_upvotes tu WHERE tu.user_id = ? AND tu.deleted_at IS NULL) AS t4 ON t.id = t4.thread_id LEFT JOIN (SELECT user_id FROM user_followers uf WHERE uf.follower_id= ?) AS t5 ON t5.user_id = t.user_id LEFT JOIN (SELECT id, thread_id FROM thread_downvotes td WHERE td.user_id = ? AND td.deleted_at IS NULL) AS t6 ON t6.thread_id = t.id WHERE t.deleted_at IS NULL;", userID, userID, userID, userID).Scan(&threads)

	if res.Error != nil {
		fmt.Println(res.Error)
		return []entity.ThreadWithDetails{}, utils.ErrInternalServerError
	}

	return threads, nil
}

func (tr *ThreadRepositoryImpl) GetThreadsFromFollowedUsers(userID uint) ([]entity.ThreadWithDetails, error) {
	var threads []entity.ThreadWithDetails

	res := tr.db.Raw("SELECT t.*, t2.upvotes_count, NOT ISNULL(t4.id) AS is_upvoted, NOT ISNULL(t3.user_id) AS is_followed, NOT ISNULL(t5.id) AS is_downvoted, users.name, users.profile_image_url, users.profession  FROM threads t LEFT JOIN (SELECT thread_id, COUNT(*) AS upvotes_count FROM thread_upvotes tu WHERE tu.deleted_at IS NULL GROUP BY thread_id) AS t2 ON t.id = t2.thread_id INNER JOIN (SELECT user_id FROM user_followers uf WHERE uf.follower_id= ?) AS t3 ON t3.user_id = t.user_id LEFT JOIN users ON users.id = t.user_id LEFT JOIN (SELECT * FROM thread_upvotes tu WHERE tu.user_id = ? AND tu.deleted_at IS NULL) AS t4 ON t.id = t4.thread_id LEFT JOIN (SELECT id, thread_id FROM thread_downvotes td WHERE td.user_id = ? AND td.deleted_at IS NULL) AS t5 ON t5.thread_id = t.id;", userID, userID, userID).Scan(&threads)

	if res.Error != nil {
		fmt.Println(res.Error)
		return []entity.ThreadWithDetails{}, utils.ErrInternalServerError
	}

	return threads, nil
}

func (tr *ThreadRepositoryImpl) AddThreadComment(comment entity.Comment) error {
	res := tr.db.Create(&comment)

	if res.Error != nil {
		fmt.Println(res.Error)
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetCommentsByThreadID(threadID uint) ([]entity.CommentDetails, error) {
	var comments []entity.CommentDetails
	res := tr.db.Raw("SELECT comments.*, users.*, t2.likes_count FROM comments LEFT JOIN (SELECT comment_id, COUNT(*) AS likes_count FROM comment_likes cl GROUP BY comment_id) AS t2 ON comments.id = t2.comment_id INNER JOIN users ON comments.user_id = users.id WHERE comments.thread_id = ? AND comments.deleted_at IS NULL", threadID).Scan(&comments)

	if res.Error != nil {
		fmt.Println(res.Error)
		return []entity.CommentDetails{}, utils.ErrInternalServerError
	}

	return comments, nil
}

func (tr *ThreadRepositoryImpl) GetThreads(keyword string, userID uint) ([]entity.ThreadWithDetails, error) {
	var threads []entity.ThreadWithDetails

	res := tr.db.Raw("SELECT combined.*, upvotes_count, NOT ISNULL(t4.id) AS is_upvoted, NOT ISNULL(t3.user_id) AS is_followed, NOT ISNULL(t5.id) AS is_downvoted, users.name, users.profile_image_url, users.profession FROM (SELECT * FROM threads t WHERE (t.body LIKE ? OR t.title LIKE ?) AND (t.deleted_at IS NULL) UNION SELECT t.* FROM comments c LEFT JOIN threads t ON t.id = c.thread_id WHERE c.body LIKE ? AND c.deleted_at IS NULL) AS combined LEFT JOIN (SELECT thread_id, COUNT(*) AS upvotes_count FROM thread_upvotes tu WHERE tu.deleted_at IS NULL GROUP BY thread_id) AS t2 ON combined.id = t2.thread_id LEFT JOIN (SELECT user_id FROM user_followers uf WHERE uf.follower_id= ?) AS t3 ON t3.user_id = combined.user_id LEFT JOIN users ON users.id = combined.user_id LEFT JOIN (SELECT * FROM thread_upvotes tu WHERE tu.user_id = ? AND tu.deleted_at IS NULL) AS t4 ON combined.id = t4.thread_id LEFT JOIN (SELECT id, thread_id, user_id FROM thread_downvotes td WHERE td.user_id = ? AND td.deleted_at IS NULL) AS t5 ON t5.thread_id = combined.id;", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", userID, userID, userID).Scan(&threads)

	if res.Error != nil {
		fmt.Println(res.Error)
		return []entity.ThreadWithDetails{}, utils.ErrInternalServerError
	}

	return threads, nil
}

func (tr *ThreadRepositoryImpl) LikeComment(commentLikes entity.CommentLikes) error {
	res := tr.db.Create(&commentLikes)

	if res.Error != nil {
		fmt.Println(res.Error)
		if strings.HasPrefix(res.Error.Error(), "Error 1452: Cannot add or update a child row") {
			return utils.ErrNotFound
		}
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) DownvoteThread(downvote entity.ThreadDownvote) error {
	res := tr.db.Create(&downvote)

	if res.Error != nil {
		fmt.Println(res.Error)
		if strings.HasPrefix(res.Error.Error(), "Error 1452: Cannot add or update a child row") {
			return utils.ErrNotFound
		} else if strings.HasPrefix(res.Error.Error(), "Error 1062: Duplicate entry") {
			return utils.ErrDuplicateEntry
		}
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) UndoDownvoteThread(threadID, userID uint) error {
	res := tr.db.Unscoped().Delete(&entity.ThreadDownvote{}, "thread_id = ? AND user_id = ?", threadID, userID)

	if res.Error != nil {
		fmt.Println(res.Error)
		return utils.ErrInternalServerError
	}

	if res.RowsAffected < 1 {
		return utils.ErrNotFound
	}

	return nil
}

func (tr *ThreadRepositoryImpl) UnlikeComment(commentID, userID uint) error {
	res := tr.db.Unscoped().Delete(&entity.CommentLikes{}, "thread_id = ? AND user_id = ?", commentID, userID)

	if res.Error != nil {
		fmt.Println(res.Error)
		if res.Error.Error() == "record not found" {
			return utils.ErrNotFound
		}
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) UndoUpvoteThread(commentID, userID uint) error {
	res := tr.db.Unscoped().Delete(&entity.ThreadUpvote{}, "thread_id = ? AND user_id = ?", commentID, userID)

	if res.Error != nil {
		fmt.Println(res.Error)
		return utils.ErrInternalServerError
	}

	if res.RowsAffected < 1 {
		return utils.ErrNotFound
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetThreadDownvotes(threadID, userID uint) (entity.ThreadDownvote, error) {
	var threadDownvotes entity.ThreadDownvote
	res := tr.db.First(&threadDownvotes, "thread_id = ? AND user_id = ?", threadID, userID)

	if res.Error != nil {
		fmt.Println(res.Error)
		if res.Error.Error() == "record not found" {
			return threadDownvotes, utils.ErrNotFound
		}
		return threadDownvotes, utils.ErrInternalServerError
	}

	return threadDownvotes, nil
}

func (tr *ThreadRepositoryImpl) GetThreadUpvotes(threadID, userID uint) (entity.ThreadUpvote, error) {
	var threadUpvote entity.ThreadUpvote
	res := tr.db.First(&threadUpvote, "thread_id = ? AND user_id = ?", threadID, userID)

	if res.Error != nil {
		fmt.Println(res.Error)
		if res.Error.Error() == "record not found" {
			return threadUpvote, utils.ErrNotFound
		}

		return threadUpvote, utils.ErrInternalServerError
	}

	return threadUpvote, nil
}

func (tr *ThreadRepositoryImpl) DeleteComment(commentID uint) error {
	res := tr.db.Debug().Delete(&entity.Comment{}, commentID)

	if res.Error != nil {
		return utils.ErrInternalServerError
	}

	if res.RowsAffected < 1 {
		return utils.ErrNotFound
	}

	return nil
}

func (tr *ThreadRepositoryImpl) CreateThreadReport(threadReport entity.ThreadReport) error {
	res := tr.db.Create(&threadReport)

	if res.Error != nil {
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetThreadReport(id uint) (entity.ThreadReport, error) {
	threadReport := entity.ThreadReport{}

	res := tr.db.Find(&threadReport, id)
	err := res.Error
	if err != nil {
		return entity.ThreadReport{}, nil
	}

	return threadReport, nil
}
func (tr *ThreadRepositoryImpl) UpdateThreadReport(threadReport entity.ThreadReport, userID uint) error {
	res := tr.db.Model(&threadReport).Update("user_id", userID)
	err := res.Error
	if err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetCommentByID(commentID uint) (entity.Comment, error) {
	var comment entity.Comment
	res := tr.db.First(&comment, commentID)

	if res.Error != nil {
		fmt.Println(res.Error)
		if res.Error.Error() == "record not found" {
			return comment, utils.ErrNotFound
		}

		return comment, utils.ErrInternalServerError
	}

	return comment, nil
}

func (tr *ThreadRepositoryImpl) CreateCommentReport(commentReport entity.CommentReport) error {
	res := tr.db.Create(&commentReport)

	if res.Error != nil {
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetCommentReport(id uint) (entity.CommentReport, error) {
	commentReport := entity.CommentReport{}

	res := tr.db.Find(&commentReport, id)
	err := res.Error
	if err != nil {
		return entity.CommentReport{}, nil
	}

	return commentReport, nil
}

func (tr *ThreadRepositoryImpl) UpdateCommentReport(commentReport entity.CommentReport, userID uint) error {
	res := tr.db.Model(&commentReport).Update("user_id", userID)
	err := res.Error
	if err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetThreadsByUserID(userID, tokenUserID uint) ([]entity.ThreadWithDetails, error) {
	threads := []entity.ThreadWithDetails{}

	res := tr.db.Raw("SELECT t.*, tlc.count AS upvotes_count, !isnull(tl.user_id) AS is_upvoted, u.*, (u.id = ?) AS is_mine FROM threads AS t LEFT JOIN (SELECT t.thread_id, COUNT(*) AS count FROM thread_upvotes AS t GROUP BY t.thread_id) AS tlc ON t.id = tlc.thread_id LEFT JOIN (SELECT * FROM thread_upvotes WHERE user_id = ?) AS tl ON tl.thread_id = t.id LEFT JOIN (SELECT u.*, !ISNULL(uf.user_id) AS is_followed FROM `users` AS u LEFT JOIN (SELECT * FROM user_followers WHERE follower_id = ?) AS uf ON u.id = uf.user_id WHERE u.deleted_at IS NULL) AS u ON u.id = t.user_id WHERE t.user_id = ? AND t.deleted_at IS NULL", userID, userID, userID, tokenUserID).Scan(&threads)

	if res.Error != nil {
		return []entity.ThreadWithDetails{}, utils.ErrInternalServerError
	}

	return threads, nil
}

func (tr *ThreadRepositoryImpl) StoreSavedThread(savedThread entity.SavedThread) error {
	res := tr.db.Create(&savedThread)

	if res.Error != nil {
		return utils.ErrInternalServerError
	}

	return nil
}

func (tr *ThreadRepositoryImpl) GetSavedThread(userID uint) ([]entity.ThreadWithDetails, error) {
	var threads []entity.ThreadWithDetails

	res := tr.db.Raw("SELECT t.*, t2.upvotes_count, NOT ISNULL(t4.id) AS is_upvoted, NOT ISNULL(t3.user_id) AS is_followed, NOT ISNULL(t5.id) AS is_downvoted, users.name, users.profile_image_url, users.profession FROM threads t INNER JOIN saved_threads st ON st.thread_id = t.id LEFT JOIN (SELECT thread_id, COUNT(*) AS upvotes_count FROM thread_upvotes tu GROUP BY thread_id) AS t2 ON t.id = t2.thread_id LEFT JOIN (SELECT user_id FROM user_followers uf WHERE uf.follower_id= ?) AS t3 ON t3.user_id = t.user_id LEFT JOIN users ON users.id = t.user_id LEFT JOIN (SELECT * FROM thread_upvotes tu WHERE tu.user_id = ?) AS t4 ON t.id = t4.thread_id LEFT JOIN (SELECT id, thread_id, user_id FROM thread_downvotes td WHERE td.user_id = ?) AS t5 ON t5.thread_id = t.id WHERE st.user_id = ? AND t.deleted_at IS NULL", userID, userID, userID, userID).Scan(&threads)

	if res.Error != nil {
		return []entity.ThreadWithDetails{}, utils.ErrInternalServerError
	}

	return threads, nil
}
