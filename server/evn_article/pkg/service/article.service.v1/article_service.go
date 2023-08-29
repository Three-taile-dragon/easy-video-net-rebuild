package article_service_v1

import (
	"context"
	consts2 "dragonsss.cn/evn_common/consts"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/model"
	article2 "dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/article/comments"
	"dragonsss.cn/evn_common/model/common"
	notice2 "dragonsss.cn/evn_common/model/user/notice"
	"dragonsss.cn/evn_grpc/article"
	"dragonsss.com/evn_article/config"
	"dragonsss.com/evn_article/internal/dao"
	"dragonsss.com/evn_article/internal/dao/mysql"
	"dragonsss.com/evn_article/internal/database"
	"dragonsss.com/evn_article/internal/database/tran"
	"dragonsss.com/evn_article/internal/repo"
	model2 "dragonsss.com/evn_article/pkg/model"
	consts "dragonsss.com/evn_ws/utils"
	"dragonsss.com/evn_ws/utils/notice"
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"go.uber.org/zap"
	"strconv"
)

// ArticleService grpc 登陆服务 实现
type ArticleService struct {
	article.UnimplementedArticleServiceServer
	cache       repo.Cache
	transaction tran.Transaction
	articleRepo repo.ArticleRepo
}

func New() *ArticleService {
	return &ArticleService{
		cache:       dao.Rc,
		transaction: dao.NewTransaction(),
		articleRepo: mysql.NewArticleDao(),
	}
}

func (as *ArticleService) GetArticleContributionList(ctx context.Context, req *article.CommonPageInfo) (*article.CommonDataResponse, error) {
	c := context.Background()
	//获取文章
	list, err := as.articleRepo.GetArticleList(c, req)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionList GetList DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetArticleContributionListResponse(list, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionList GetArticleContributionListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &article.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (as *ArticleService) GetArticleContributionListByUser(ctx context.Context, req *article.CommonIDRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//获取文章
	list, err := as.articleRepo.GetArticleListByID(c, req)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionListByUser GetArticleListByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetArticleContributionListByUserResponse(list, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionListByUser GetArticleContributionListByUserResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionListByUser rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &article.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (as *ArticleService) GetArticleComment(ctx context.Context, req *article.GetArticleCommentRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//获取文章
	com, err := as.articleRepo.GetArticleComments(c, req)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleComment GetArticleComments DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetArticleContributionCommentsResponse(com, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleComment GetArticleContributionCommentsResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleComment rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &article.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (as *ArticleService) GetArticleClassificationList(ctx context.Context, req *article.CommonZeroRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//获取文章分类
	cn, err := as.articleRepo.FindAllClassification(c, req)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleClassificationList FindAllClassification DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp := model2.GetArticleClassificationListResponse(cn)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleClassificationList GetArticleClassificationListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleClassificationList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &article.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (as *ArticleService) GetArticleTotalInfo(ctx context.Context, req *article.CommonZeroRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//查询文章数量
	articleNm, err := as.articleRepo.GetAllArticleCount(c)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleTotalInfo GetAllArticleCount DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//获取文章分类
	cn, err := as.articleRepo.FindAllClassification(c, req)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleTotalInfo FindAllClassification DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	cnNum := int64(len(*cn))
	rsp := model2.GetArticleTotalInfoResponse(cn, articleNm, cnNum)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleTotalInfo GetArticleTotalInfoResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleTotalInfo rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &article.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (as *ArticleService) GetArticleContributionByID(ctx context.Context, req *article.CommonIDAndUIDRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//查询文章数量
	ac, err := as.articleRepo.GetArticleInfoByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionByID GetArticleInfoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if req.UID > 0 {
		//添加历史记录
		if is, err := as.articleRepo.AddArticleRecord(c, req); !is || err != nil {
			zap.L().Error("evn_article article_service GetArticleContributionByID AddArticleRecord DB_error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		//进行文章热度增加
		if !dao.Rc.R().SIsMember(c, consts.ArticleWatchByID+strconv.Itoa(int(req.ID)), strconv.Itoa(int(req.UID))).Val() {
			//最近无播放
			dao.Rc.R().SAdd(c, consts.ArticleWatchByID+strconv.Itoa(int(req.ID)), req.UID)
			//添加浏览量
			if as.articleRepo.WatchArticle(c, req.ID) != nil {
				zap.L().Error("evn_article article_service GetArticleContributionByID WatchArticle DB_error", zap.Error(err))
				return nil, errs.GrpcError(model.DBError)
			}
			ac.Heat++
		}
	}
	rsp := model2.GetArticleContributionByIDResponse(ac, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionByID GetArticleContributionByIDResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleContributionByID rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &article.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (as *ArticleService) CreateArticleContribution(ctx context.Context, req *article.CreateArticleContributionRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//进行内容判断
	for _, v := range req.Label {
		vRune := []rune(v) //避免中文占位问题
		if len(vRune) > 7 {
			return nil, fmt.Errorf("标签长度不能大于7位")
		}
	}

	coverImg, _ := json.Marshal(common.Img{
		Src: req.Cover,
		Tp:  req.CoverUploadType,
	})
	//正则匹配替换url
	//取url前缀
	prefix, err := conversion.SwitchTypeAsUrlPrefix(req.ArticleContributionUploadType, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		return nil, fmt.Errorf("保存资源方式不存在")
	}
	//正则匹配替换
	reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+prefix, 0)
	match, err := reg.Replace(req.Content, consts2.UrlPrefixSubstitution, -1, -1)
	req.Content = match
	//插入数据
	articlesContribution := article2.ArticlesContribution{
		Uid:                uint(req.Uid),
		ClassificationID:   uint(req.ClassificationID),
		Title:              req.Title,
		Cover:              coverImg,
		Label:              conversion.MapConversionString(req.Label),
		Content:            req.Content,
		ContentStorageType: req.ArticleContributionUploadType,
		IsComments:         conversion.BoolTurnInt8(req.Comments),
		Heat:               0,
	}
	//保存文章
	err = as.articleRepo.CreateArticle(c, &articlesContribution)
	if err != nil {
		zap.L().Error("evn_article article_service CreateArticleContribution CreateArticle DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	return &article.CommonDataResponse{Data: "保存成功"}, nil
}

func (as *ArticleService) UpdateArticleContribution(ctx context.Context, req *article.UpdateArticleContributionRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//更新专栏
	ac, err := as.articleRepo.GetArticleInfoByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_article article_service UpdateArticleContribution GetArticleInfoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if ac.Uid != uint(req.Uid) {
		return &article.CommonDataResponse{Data: "非法操作"}, errs.GrpcError(model.RequestError)
	}
	coverImg, _ := json.Marshal(common.Img{
		Src: req.Cover,
		Tp:  req.CoverUploadType,
	})
	updateList := map[string]interface{}{
		"cover":             coverImg,
		"title":             req.Title,
		"label":             conversion.MapConversionString(req.Label),
		"content":           req.Content,
		"is_comments":       req.Comments,
		"classification_id": req.ClassificationID,
	}
	err = as.articleRepo.UpdateArticle(c, &updateList)
	if err != nil {
		zap.L().Error("evn_article article_service UpdateArticleContribution UpdateArticle DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	return &article.CommonDataResponse{Data: "保存成功"}, nil
}

func (as *ArticleService) DeleteArticleByID(ctx context.Context, req *article.CommonIDAndUIDRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//删除专栏
	if is, err := as.articleRepo.DeleteArticleByID(c, req); !is || err != nil {
		zap.L().Error("evn_article article_service DeleteArticleByID DeleteArticleByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	return &article.CommonDataResponse{Data: "删除成功"}, nil
}

func (as *ArticleService) ArticlePostComment(ctx context.Context, req *article.ArticlePostCommentRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//查询专栏
	ac, err := as.articleRepo.GetArticleInfoByID(c, req.ArticleID)
	if err != nil {
		zap.L().Error("evn_article article_service ArticlePostComment GetArticleInfoByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	CommentFirstID, err := as.articleRepo.GetCommentFirstIDByID(c, req.ContentID)
	if err != nil {
		zap.L().Error("evn_article article_service ArticlePostComment GetCommentFirstIDByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	CommentUserID, err := as.articleRepo.GetCommentUserIDByID(c, req.ContentID)
	if err != nil {
		zap.L().Error("evn_article article_service ArticlePostComment GetCommentUserIDByID DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	comment := comments.Comment{
		Uid:            uint(req.Uid),
		ArticleID:      uint(req.ArticleID),
		Context:        req.Content,
		CommentID:      uint(req.ContentID),
		CommentUserID:  CommentUserID.Uid,
		CommentFirstID: CommentFirstID.ID,
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = as.transaction.Action(func(conn database.DbConn) error {
		err = as.articleRepo.CreateArticleComment(conn, c, &comment)
		if err != nil {
			zap.L().Error("evn_video video_service VideoPostComment CreateComment Tx_DB_error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		//消息通知
		if ac.Uid == comment.Uid {
			return nil
		}
		//添加消息通知
		err = as.articleRepo.AddNotice(c, ac.Uid, comment.Uid, ac.ID, notice2.ArticleComment, comment.Context)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &article.CommonDataResponse{Data: "发布失败"}, err
	}
	//socket推送(在线的情况下)
	if _, ok := notice.Severe.UserMapChannel[ac.UserInfo.ID]; ok {
		userChannel := notice.Severe.UserMapChannel[ac.UserInfo.ID]
		userChannel.NoticeMessage(notice2.ArticleComment)
	}
	return &article.CommonDataResponse{Data: "删除成功"}, nil
}

func (as *ArticleService) GetArticleManagementList(ctx context.Context, req *article.GetArticleManagementListRequest) (*article.CommonDataResponse, error) {
	c := context.Background()
	//获取个人发布专栏信息
	list, err := as.articleRepo.GetArticleManagementList(c, req)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleManagementList GetArticleManagementList DB_error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp, err := model2.GetArticleManagementListResponse(list, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleManagementList GetArticleManagementListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_article article_service GetArticleManagementList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &article.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}
