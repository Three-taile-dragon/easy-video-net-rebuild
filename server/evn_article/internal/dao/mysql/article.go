package mysql

import (
	"context"
	"dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/article/classification"
	comments2 "dragonsss.cn/evn_common/model/article/comments"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/notice"
	"dragonsss.cn/evn_common/model/user/record"
	article2 "dragonsss.cn/evn_grpc/article"
	"dragonsss.com/evn_article/internal/database"
	"dragonsss.com/evn_article/internal/database/gorms"
	"gorm.io/gorm"
)

type ArticleDao struct {
	conn *gorms.GormConn
}

func NewArticleDao() *ArticleDao {
	return &ArticleDao{
		conn: gorms.New(),
	}
}

func (a ArticleDao) GetUserByID(ctx context.Context, uid uint32) (*user.User, error) {
	session := a.conn.Session(ctx)
	var user1 *user.User
	err := session.
		Where("id", uid).
		Find(&user1).
		Error
	if err != nil {
		return nil, err
	}
	return user1, nil
}

func (a ArticleDao) GetArticleList(ctx context.Context, req *article2.CommonPageInfo) (*article.ArticlesContributionList, error) {
	session := a.conn.Session(ctx)
	var ar *article.ArticlesContributionList
	err := session.
		Preload("Likes").
		Preload("Classification").
		Preload("UserInfo").
		Preload("Comments").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size)).
		Order("created_at desc").
		Find(&ar).
		Error
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (a ArticleDao) GetArticleListByID(ctx context.Context, req *article2.CommonIDRequest) (*article.ArticlesContributionList, error) {
	session := a.conn.Session(ctx)
	var ar *article.ArticlesContributionList
	err := session.
		Where("uid", req.ID).
		Preload("Likes").
		Preload("Classification").
		Preload("Comments").
		Order("created_at desc").
		Find(&ar).
		Error
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (a ArticleDao) GetArticleComments(ctx context.Context, req *article2.GetArticleCommentRequest) (*article.ArticlesContribution, error) {
	session := a.conn.Session(ctx)
	var ac *article.ArticlesContribution
	err := session.
		Where("id", req.ArticleID).
		Preload("Likes").
		Preload("Classification").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").Order("created_at desc").Limit(int(req.PageInfo.Size)).Offset(int((req.PageInfo.Page - 1) * req.PageInfo.Size))
		}).
		Find(&ac).
		Error
	if err != nil {
		return nil, err
	}
	return ac, nil
}

func (a ArticleDao) FindAllClassification(ctx context.Context, req *article2.CommonZeroRequest) (*classification.ClassificationsList, error) {
	session := a.conn.Session(ctx)
	var cn *classification.ClassificationsList
	err := session.
		Find(&cn).
		Error
	if err != nil {
		return nil, err
	}
	return cn, nil
}

func (a ArticleDao) GetAllArticleCount(ctx context.Context) (*int64, error) {
	session := a.conn.Session(ctx)
	num := new(int64)
	var ac *article.ArticlesContributionList
	err := session.
		Find(&ac).
		Count(num).
		Error
	if err != nil {
		return nil, err
	}
	return num, nil
}

func (a ArticleDao) GetArticleInfoByID(ctx context.Context, id uint32) (*article.ArticlesContribution, error) {
	session := a.conn.Session(ctx)
	var ac *article.ArticlesContribution
	err := session.
		Where("id", id).
		Preload("Likes").
		Preload("UserInfo").
		Preload("Classification").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserInfo").Order("created_at desc")
		}).
		Find(&ac).
		Error
	if err != nil {
		return nil, err
	}
	return ac, nil
}

func (a ArticleDao) AddArticleRecord(ctx context.Context, req *article2.CommonIDAndUIDRequest) (bool, error) {
	session := a.conn.Session(ctx)
	var rc *record.Record
	err := session.
		Where(record.Record{Uid: uint(req.UID), Type: "article", ToId: uint(req.ID)}).
		Find(&rc).
		Error
	if err != nil {
		return false, err
	}
	if rc.ID <= 0 {
		//创建记录
		rc.Uid = uint(req.UID)
		rc.Type = "article"
		rc.ToId = uint(req.ID)
		err := session.
			Create(&rc).
			Error
		if err != nil {
			return false, err
		}
	} else {
		//更新记录
		err := session.
			Where("id", rc.ID).
			Updates(&rc).
			Error
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (a ArticleDao) WatchArticle(ctx context.Context, id uint32) error {
	session := a.conn.Session(ctx)
	var ac *article.ArticlesContribution
	err := session.
		Model(ac).
		Where("id", id).
		Updates(map[string]interface{}{"heat": gorm.Expr("Heat  + ?", 1)}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleDao) CreateArticle(ctx context.Context, articlesContribution *article.ArticlesContribution) error {
	session := a.conn.Session(ctx)
	err := session.
		Create(articlesContribution).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleDao) UpdateArticle(ctx context.Context, m *map[string]interface{}) error {
	session := a.conn.Session(ctx)
	err := session.
		Model(&article.ArticlesContribution{}).
		Create(m).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleDao) DeleteArticleByID(ctx context.Context, req *article2.CommonIDAndUIDRequest) (bool, error) {
	session := a.conn.Session(ctx)
	var ac *article.ArticlesContribution
	err := session.
		Where("id", req.ID).
		Find(&ac).
		Error
	if err != nil {
		return false, err
	}
	if ac.Uid != uint(req.UID) {
		return false, nil
	}
	err = session.
		Delete(&ac).
		Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a ArticleDao) GetCommentFirstIDByID(ctx context.Context, contentID uint32) (*comments2.Comment, error) {
	session := a.conn.Session(ctx)
	var com *comments2.Comment
	err := session.
		Where("id", contentID).
		Find(&com).
		Error
	if err != nil {
		return nil, err
	}
	//循环获取最顶层的评论ID
	for com.CommentID != 0 {
		c := context.Background()
		com, err = a.GetCommentFirstIDByID(c, uint32(com.CommentID))
		if err != nil {
			return nil, err
		}
	}
	return com, nil
}

func (a ArticleDao) GetCommentUserIDByID(ctx context.Context, contentID uint32) (*comments2.Comment, error) {
	session := a.conn.Session(ctx)
	var com *comments2.Comment
	err := session.
		Where("id", contentID).
		Find(&com).
		Error
	if err != nil {
		return nil, err
	}
	return com, nil
}

func (a ArticleDao) CreateArticleComment(conn database.DbConn, ctx context.Context, c *comments2.Comment) error {
	a.conn = conn.(*gorms.GormConn) //使用事务操作
	return a.conn.Tx(ctx).Create(c).Error
}

func (a ArticleDao) AddNotice(ctx context.Context, uid uint, cid uint, tid uint, tp string, c string) error {
	session := a.conn.Session(ctx)
	err := session.
		Create(&notice.Notice{
			Uid:     uid,
			Cid:     cid,
			ToID:    tid,
			Type:    tp,
			Content: c,
			ISRead:  0,
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (a ArticleDao) GetArticleManagementList(ctx context.Context, req *article2.GetArticleManagementListRequest) (*article.ArticlesContributionList, error) {
	session := a.conn.Session(ctx)
	var ac *article.ArticlesContributionList
	err := session.
		Where("uid", req.Uid).
		Preload("Likes").
		Preload("Classification").
		Preload("Comments").
		Limit(int(req.PageInfo.Size)).Offset(int((req.PageInfo.Page - 1) * req.PageInfo.Size)).
		Order("created_at desc").
		Find(&ac).
		Error
	if err != nil {
		return nil, err
	}
	return ac, nil
}
