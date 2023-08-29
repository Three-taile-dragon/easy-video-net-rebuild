package model

import (
	"context"
	"dragonsss.cn/evn_common/consts"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/model/article"
	"dragonsss.cn/evn_common/model/article/classification"
	"dragonsss.cn/evn_common/model/common"
	"dragonsss.com/evn_article/internal/dao/mysql"
	"encoding/json"
	"github.com/dlclark/regexp2"
	"time"
)

type GetArticleContributionListByUserResponseStruct struct {
	Id             uint      `json:"id"`
	Uid            uint      `json:"uid" `
	Title          string    `json:"title" `
	Cover          string    `json:"cover" `
	Label          []string  `json:"label" `
	Content        string    `json:"content"`
	IsComments     bool      `json:"is_comments"`
	Heat           int       `json:"heat"`
	LikesNumber    int       `json:"likes_number"`
	CommentsNumber int       `json:"comments_number"`
	Classification string    `json:"classification"`
	CreatedAt      time.Time `json:"created_at"`
}

// 评论信息
type commentsInfo struct {
	ID              uint             `json:"id"`
	CommentID       uint             `json:"comment_id"`
	CommentFirstID  uint             `json:"comment_first_id"`
	CreatedAt       time.Time        `json:"created_at"`
	Context         string           `json:"context"`
	Uid             uint             `json:"uid"`
	Username        string           `json:"username"`
	Photo           string           `json:"photo"`
	CommentUserID   uint             `json:"comment_user_id"`
	CommentUserName string           `json:"comment_user_name"`
	LowerComments   commentsInfoList `json:"lowerComments"`
}

type commentsInfoList []*commentsInfo

type GetArticleContributionCommentsResponseStruct struct {
	Id             uint             `json:"id"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
}

type GetArticleContributionListByUserResponseList []GetArticleContributionListByUserResponseStruct

// ArticleClassificationInfo 文章分类信息
type ArticleClassificationInfo struct {
	ID       uint                          `json:"id"`
	AID      uint                          `json:"aid"`
	Label    string                        `json:"label"`
	Children ArticleClassificationInfoList `json:"children"`
}

type ArticleClassificationInfoList []*ArticleClassificationInfo

type GetArticleTotalInfoResponseStruct struct {
	Classification    ArticleClassificationInfoList `json:"classification"`
	ArticleNum        int64                         `json:"article_num"`
	ClassificationNum int64                         `json:"classification_num"`
}

type GetArticleContributionByIDResponseStruct struct {
	Id             uint             `json:"id"`
	Uid            uint             `json:"uid" `
	Title          string           `json:"title" `
	Cover          string           `json:"cover" `
	Label          []string         `json:"label" `
	Content        string           `json:"content"`
	IsComments     bool             `json:"is_comments"`
	Heat           int              `json:"heat"`
	LikesNumber    int              `json:"likes_number"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
	CreatedAt      time.Time        `json:"created_at"`
}

type GetArticleManagementListItem struct {
	ID               uint     `json:"id"`
	ClassificationID uint     `json:"classification_id"`
	Title            string   `json:"title"`
	Cover            string   `json:"cover"`
	CoverUrl         string   `json:"cover_url"`
	CoverUploadType  string   `json:"cover_upload_type"`
	Label            []string `json:"label"`
	Content          string   `json:"content"`
	IsComments       bool     `json:"is_comments" `
	Heat             int      `json:"heat"`
}

type GetArticleManagementListResponseStruct []GetArticleManagementListItem

func GetArticleContributionListResponse(l *article.ArticlesContributionList, localhost string, tencentOssHost string) GetArticleContributionListByUserResponseList {
	response := make(GetArticleContributionListByUserResponseList, 0)
	for _, v := range *l {
		coverSrc, _ := conversion.FormattingJsonSrc(v.Cover, localhost, tencentOssHost)

		//正则替换首文内容
		reg := regexp2.MustCompile(`<(\S*?)[^>]*>.*?|<.*? />`, 0)
		match, _ := reg.Replace(v.Content, "", -1, -1)
		matchRune := []rune(match)
		if len(matchRune) > 100 {
			v.Content = string(matchRune[:100]) + "..."
		} else {
			v.Content = match
		}

		//只显示一个标签
		label := conversion.StringConversionMap(v.Label)
		if len(label) >= 2 {
			label = label[:1]
		}

		response = append(response, GetArticleContributionListByUserResponseStruct{
			Id:             v.ID,
			Uid:            v.Uid,
			Title:          v.Title,
			Cover:          coverSrc,
			Label:          label,
			Content:        v.Content,
			Classification: v.Classification.Label,
			IsComments:     conversion.Int8TurnBool(v.IsComments),
			Heat:           v.Heat,
			LikesNumber:    len(v.Likes),
			CommentsNumber: len(v.Comments),
			CreatedAt:      v.CreatedAt,
		})
	}
	return response
}

func GetArticleContributionListByUserResponse(l *article.ArticlesContributionList, localhost string, tencentOssHost string) GetArticleContributionListByUserResponseList {
	response := make(GetArticleContributionListByUserResponseList, 0)
	for _, v := range *l {
		coverSrc, _ := conversion.FormattingJsonSrc(v.Cover, localhost, tencentOssHost)

		//正则替换首文内容
		reg := regexp2.MustCompile(`<(\S*?)[^>]*>.*?|<.*? />`, 0)
		match, _ := reg.Replace(v.Content, "", -1, -1)
		matchRune := []rune(match)
		if len(matchRune) > 100 {
			v.Content = string(matchRune[:100]) + "..."
		} else {
			v.Content = match
		}

		//只显示一个标签
		label := conversion.StringConversionMap(v.Label)
		if len(label) >= 3 {
			label = label[:1]
		}

		response = append(response, GetArticleContributionListByUserResponseStruct{
			Id:             v.ID,
			Uid:            v.Uid,
			Title:          v.Title,
			Cover:          coverSrc,
			Label:          label,
			Content:        v.Content,
			Classification: v.Classification.Label,
			IsComments:     conversion.Int8TurnBool(v.IsComments),
			Heat:           v.Heat,
			LikesNumber:    len(v.Likes),
			CommentsNumber: len(v.Comments),
			CreatedAt:      v.CreatedAt,
		})
	}
	return response
}

func GetArticleContributionCommentsResponse(vc *article.ArticlesContribution, localhost string, tencentOssHost string) GetArticleContributionCommentsResponseStruct {
	//评论
	comments := commentsInfoList{}
	for _, v := range vc.Comments {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo, localhost, tencentOssHost)
		articleRepo := mysql.NewArticleDao()
		ctx := context.Background()
		commentUser, _ := articleRepo.GetUserByID(ctx, uint32(v.CommentUserID))
		comments = append(comments, &commentsInfo{
			ID:              v.ID,
			CommentID:       v.CommentID,
			CommentFirstID:  v.CommentFirstID,
			CommentUserID:   v.CommentUserID,
			CommentUserName: commentUser.Username,
			CreatedAt:       v.CreatedAt,
			Context:         v.Context,
			Uid:             v.UserInfo.ID,
			Username:        v.UserInfo.Username,
			Photo:           photo,
		})
	}
	commentsList := comments.getChildComment()
	//输出
	response := GetArticleContributionCommentsResponseStruct{
		Id:             vc.ID,
		Comments:       commentsList,
		CommentsNumber: len(vc.Comments),
	}
	return response
}

// 得到分级结构
func (l commentsInfoList) getChildComment() commentsInfoList {
	topList := commentsInfoList{}
	for _, v := range l {
		if v.CommentID == 0 {
			//顶层
			topList = append(topList, v)
		}
	}
	return commentsInfoListSecondTree(topList, l)
}

func commentsInfoListSecondTree(menus commentsInfoList, allData commentsInfoList) commentsInfoList {
	//循环所有一级菜单
	for k, v := range menus {
		//查询所有该菜单下的所有子菜单
		var nodes commentsInfoList //定义子节点目录
		for _, av := range allData {
			if av.CommentFirstID == v.ID {
				nodes = append(nodes, av)
			}
		}
		for kk, _ := range nodes {
			menus[k].LowerComments = append(menus[k].LowerComments, nodes[kk])
		}
		//将刚刚查询出来的子菜单进行递归,查询出三级菜单和四级菜单
		commentsInfoListTree(nodes, allData)
	}
	return menus
}

// 生成树结构
func commentsInfoListTree(menus commentsInfoList, allData commentsInfoList) commentsInfoList {
	//循环所有一级菜单
	for k, v := range menus {
		//查询所有该菜单下的所有子菜单
		var nodes commentsInfoList //定义子节点目录
		for _, av := range allData {
			if av.CommentID == v.ID {
				nodes = append(nodes, av)
			}
		}
		for kk, _ := range nodes {
			menus[k].LowerComments = append(menus[k].LowerComments, nodes[kk])
		}
		//将刚刚查询出来的子菜单进行递归,查询出三级菜单和四级菜单
		commentsInfoListTree(nodes, allData)
	}
	return menus
}

func GetArticleClassificationListResponse(cl *classification.ClassificationsList) ArticleClassificationInfoList {
	response := make(ArticleClassificationInfoList, 0)
	for _, v := range *cl {
		response = append(response, &ArticleClassificationInfo{
			ID:       v.ID,
			AID:      v.AID,
			Label:    v.Label,
			Children: make(ArticleClassificationInfoList, 0),
		})
	}
	return response.getChildComment()
}

// 得到分级结构
func (l ArticleClassificationInfoList) getChildComment() ArticleClassificationInfoList {
	topList := ArticleClassificationInfoList{}
	for _, v := range l {
		if v.AID == 0 {
			//顶层
			topList = append(topList, &ArticleClassificationInfo{
				ID:       v.ID,
				AID:      v.AID,
				Label:    v.Label,
				Children: nil,
			})
		}
	}
	return classificationInfoListTree(topList, l)
}

// 生成树结构
func classificationInfoListTree(menus ArticleClassificationInfoList, allData ArticleClassificationInfoList) ArticleClassificationInfoList {
	//循环所有一级菜单
	for k, v := range menus {
		//查询所有该菜单下的所有子菜单
		var nodes ArticleClassificationInfoList //定义子节点目录
		for _, av := range allData {
			if av.AID == v.ID {
				nodes = append(nodes, &ArticleClassificationInfo{
					ID:       av.ID,
					AID:      av.AID,
					Label:    av.Label,
					Children: nil,
				})
			}
		}
		for kk, _ := range nodes {
			menus[k].Children = append(menus[k].Children, nodes[kk])
		}
		//将刚刚查询出来的子菜单进行递归,查询出三级菜单和四级菜单
		classificationListSecondTree(nodes, allData)
	}
	return menus
}

func classificationListSecondTree(menus ArticleClassificationInfoList, allData ArticleClassificationInfoList) ArticleClassificationInfoList {
	//循环所有一级菜单
	for k, v := range menus {
		//查询所有该菜单下的所有子菜单
		var nodes ArticleClassificationInfoList //定义子节点目录
		for _, av := range allData {
			if av.AID == v.ID {
				nodes = append(nodes, av)
			}
		}
		for kk, _ := range nodes {
			menus[k].Children = append(menus[k].Children, nodes[kk])
		}
		//将刚刚查询出来的子菜单进行递归,查询出三级菜单和四级菜单
		classificationListSecondTree(nodes, allData)
	}
	return menus
}

func GetArticleTotalInfoResponse(cl *classification.ClassificationsList, articleNum *int64, clNum int64) interface{} {
	classificationInfo := make(ArticleClassificationInfoList, 0)
	for _, v := range *cl {
		classificationInfo = append(classificationInfo, &ArticleClassificationInfo{
			ID:       v.ID,
			AID:      v.AID,
			Label:    v.Label,
			Children: make(ArticleClassificationInfoList, 0),
		})
	}
	classificationInfo = classificationInfo.getChildComment()

	return GetArticleTotalInfoResponseStruct{
		Classification:    classificationInfo,
		ArticleNum:        *articleNum,
		ClassificationNum: clNum,
	}
}

func GetArticleContributionByIDResponse(vc *article.ArticlesContribution, localhost string, tencentOssHost string) GetArticleContributionByIDResponseStruct {
	coverSrc, _ := conversion.FormattingJsonSrc(vc.Cover, localhost, tencentOssHost)

	prefix, _ := conversion.SwitchTypeAsUrlPrefix(vc.ContentStorageType, localhost, tencentOssHost)
	//正则替换src
	reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+consts.UrlPrefixSubstitutionEscape, 0)
	match, _ := reg.Replace(vc.Content, prefix, -1, -1)
	vc.Content = match

	label := conversion.StringConversionMap(vc.Label)
	if len(label) >= 2 {
		label = label[:1]
	}
	//评论
	comments := commentsInfoList{}
	for _, v := range vc.Comments {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo, localhost, tencentOssHost)
		articleRepo := mysql.NewArticleDao()
		ctx := context.Background()
		commentUser, _ := articleRepo.GetUserByID(ctx, uint32(v.CommentUserID))
		comments = append(comments, &commentsInfo{
			ID:              v.ID,
			CommentID:       v.CommentID,
			CommentFirstID:  v.CommentFirstID,
			CommentUserID:   v.CommentUserID,
			CommentUserName: commentUser.Username,
			CreatedAt:       v.CreatedAt,
			Context:         v.Context,
			Uid:             v.UserInfo.ID,
			Username:        v.UserInfo.Username,
			Photo:           photo,
		})
	}
	commentsList := comments.getChildComment()

	//输出
	response := GetArticleContributionByIDResponseStruct{
		Id:             vc.ID,
		Uid:            vc.Uid,
		Title:          vc.Title,
		Cover:          coverSrc,
		Label:          label,
		Content:        vc.Content,
		IsComments:     conversion.Int8TurnBool(vc.IsComments),
		Heat:           vc.Heat,
		LikesNumber:    len(vc.Likes),
		Comments:       commentsList,
		CommentsNumber: len(vc.Comments),
		CreatedAt:      vc.CreatedAt,
	}
	return response
}

func GetArticleManagementListResponse(al *article.ArticlesContributionList, localhost string, tencentOssHost string) (interface{}, error) {
	list := make(GetArticleManagementListResponseStruct, 0)
	for _, v := range *al {
		coverJson := new(common.Img)
		_ = json.Unmarshal(v.Cover, coverJson)
		cover, _ := conversion.FormattingJsonSrc(v.Cover, localhost, tencentOssHost)
		prefix, _ := conversion.SwitchTypeAsUrlPrefix(v.ContentStorageType, localhost, tencentOssHost)
		//正则替换src
		reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+consts.UrlPrefixSubstitutionEscape, 0)
		match, _ := reg.Replace(v.Content, prefix, -1, -1)
		v.Content = match

		list = append(list, GetArticleManagementListItem{
			ID:               v.ID,
			ClassificationID: v.ClassificationID,
			Title:            v.Title,
			Cover:            cover,
			CoverUrl:         coverJson.Src,
			CoverUploadType:  coverJson.Tp,
			Label:            conversion.StringConversionMap(v.Label),
			Content:          v.Content,
			IsComments:       conversion.Int8TurnBool(v.IsComments),
			Heat:             v.Heat,
		})
	}

	return list, nil
}
