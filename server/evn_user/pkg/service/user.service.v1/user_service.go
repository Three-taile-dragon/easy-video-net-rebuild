package user_service_v1

import (
	"context"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/conversion"
	"dragonsss.cn/evn_common/encrypts"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/jwts"
	"dragonsss.cn/evn_common/model"
	common2 "dragonsss.cn/evn_common/model/common"
	mCommon "dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/liveInfo"
	"dragonsss.cn/evn_common/model/user"
	"dragonsss.cn/evn_common/model/user/chat/chatList"
	"dragonsss.cn/evn_common/model/user/chat/chatMsg"
	"dragonsss.cn/evn_common/model/user/collect"
	"dragonsss.cn/evn_common/model/user/favorites"
	"dragonsss.cn/evn_common/model/user/notice"
	"dragonsss.cn/evn_common/response"
	user2 "dragonsss.cn/evn_grpc/user"
	"dragonsss.cn/evn_user/config"
	"dragonsss.cn/evn_user/internal/dao"
	"dragonsss.cn/evn_user/internal/dao/mysql"
	"dragonsss.cn/evn_user/internal/database"
	"dragonsss.cn/evn_user/internal/database/tran"
	"dragonsss.cn/evn_user/internal/repo"
	"dragonsss.cn/evn_user/util"
	email1 "dragonsss.cn/evn_user/util/email"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"strconv"
	"strings"
	"time"
)

// UserService grpc 登陆服务实现
type UserService struct {
	user2.UnimplementedUserServiceServer
	cache       repo.Cache
	userRepo    repo.UserRepo
	transaction tran.Transaction
	//memberRepo  repo.MemberRepo
}

func New() *UserService {
	return &UserService{
		cache:       dao.Rc,
		userRepo:    mysql.NewUserDao(),
		transaction: dao.NewTransaction(),
		//memberRepo:  mysql.NewMemberDao(),
	}
}

func (ls *UserService) GetCaptcha(ctx context.Context, req *user2.CaptchaRequest) (*user2.CommonDataResponse, error) {
	//1.获取参数
	email := req.Email
	//2.校验参数
	if !common.VerifyEmailFormat(email) {
		return nil, errs.GrpcError(model.NoLegalEmail) //使用自定义错误码进行处理
	}
	//3.生成验证码(随机四位1000-9999或者六位100000-999999)
	code := util.CreateCaptcha(6) //生成随机六位数字验证码
	fmt.Printf("%v验证码为：%v", email, code)
	//4.调用短信平台(第三方 放入go func 协程 接口可以快速响应
	//TODO 完善邮件发送服务
	go func() {
		//发送方
		mailTo := []string{req.Email}
		// 邮件主题
		subject := "验证码"
		// 邮件正文
		body := fmt.Sprintf("您正在注册验证码为:%s,5分钟有效,请勿转发他人", code)
		//TODO 测试暂时忽略错误
		_ = email1.SendMail(mailTo, subject, body)
		//redis存储	假设后续缓存可能存在mysql当中,也可以存在mongo当中,也可能存在memcache当中
		//使用接口 达到低耦合高内聚
		//5.存储验证码 redis 当中,过期时间5分钟
		//redis.Set"REGISTER_"+mobile, code)
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
		defer cancel()
		err := ls.cache.Put(c, model.RegisterRedisKey+email, code, 5*time.Minute)
		if err != nil {
			zap.L().Error("evn_user user_service GetCaptcha redis put err", zap.Error(err))

		}
		//zap.L().Debug("将手机号和验证码存入redis成功：REGISTER_" + email + " : " + code + "\n")
	}()
	//注意code一般不发送
	//这里是做了简化处理 由于短信平台目前对于个人不好使用
	return &user2.CommonDataResponse{Data: "发送成功"}, nil
}

func (ls *UserService) Register(ctx context.Context, req *user2.RegisterRequest) (*user2.UserInfoResponse, error) {
	c := context.Background()
	//可以校验参数
	//校验验证码
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+req.Email)
	if err == redis.Nil {
		return nil, errs.GrpcError(model.CaptchaNoExist)
	}
	if err != nil {
		zap.L().Error("evn_user user_service Register redis read err", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != req.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	//校验业务逻辑
	exist, err := ls.userRepo.IsExistByEmail(c, req.Email)
	if err != nil {
		zap.L().Error("evn_user user_service Register GetUserByEmail DB_Error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	//检验用户名
	exist, err = ls.userRepo.IsExistByName(c, req.Name)
	if err != nil {
		zap.L().Error("evn_user user_service Register GetUserByName DB_Error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}
	////检验手机号
	//exist, err = ls.userRepo.GetUserByMobile(c, req.Mobile)
	//if err != nil {
	//	zap.L().Error("数据库出错", zap.Error(err))
	//	return nil, errs.GrpcError(model.DBError)
	//}
	//if exist {
	//	return nil, errs.GrpcError(model.MobileExist)
	//}
	//执行业务逻辑
	//pwd := encrypts.Md5(req.Password) //加密部分
	//随机生成用户ID
	//使用薄雾算法生成user id
	//mist := common.NewMist()
	//userIdSequence := mist.Generate()

	//bcrypt 密码加密
	pwHashBytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	//转换成字符串
	pwHashStr := string(pwHashBytes)

	photo, _ := json.Marshal(mCommon.Img{
		Src: "",
		Tp:  "local",
	})

	mem := &user.User{
		Email:     req.Email,
		Username:  req.Name,
		Password:  pwHashStr,
		Photo:     photo,
		BirthDate: time.Now(),
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.userRepo.SaveUser(conn, c, mem)
		if err != nil {
			zap.L().Error("evn_user user_service Register SaveUser DB_Error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	//var conn database.DbConn
	//err = ls.userRepo.SaveMember(conn, c, mem)
	////err = ls.memberRepo.SaveMember(c, mem)
	//使用jwt生成token
	memIdStr := strconv.FormatInt(int64(mem.ID), 10)
	token := jwts.CreateToken(memIdStr, config.C.JC.AccessExp, config.C.JC.AccessSecret, config.C.JC.RefreshSecret, config.C.JC.RefreshExp)
	return &user2.UserInfoResponse{
		Id:        int64(mem.ID),
		Username:  mem.Username,
		Photo:     "",
		Token:     token.AccessToken,
		CreatedAt: mem.CreatedAt.String(),
	}, nil
}

func (ls *UserService) Login(ctx context.Context, req *user2.LoginRequest) (*user2.UserInfoResponse, error) {
	c := context.Background()
	//获取传入参数
	//校验参数
	//校验用户名和邮箱
	exist, err := ls.userRepo.IsExistByName(c, req.Username)
	if err != nil {
		zap.L().Error("evn_user user_service Login GetUserByNameAndEmail DB_Error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if !exist {
		return nil, errs.GrpcError(model.AccountNoExist)
	}
	//查询账号密码是否正确

	mem, err := ls.userRepo.CheckPassword(c, req.Username)
	//若err不为空说明密码不匹配
	err = bcrypt.CompareHashAndPassword([]byte(mem.Password), []byte(req.Password))
	if err != nil {
		//zap.L().Error("登陆模块member数据库查询出错", zap.Error(err))
		return nil, errs.GrpcError(model.AccountAndPwdError)
	}

	//使用jwt生成token
	memIdStr := strconv.FormatInt(int64(mem.ID), 10)
	token := jwts.CreateToken(memIdStr, config.C.JC.AccessExp, config.C.JC.AccessSecret, config.C.JC.RefreshSecret, config.C.JC.RefreshExp)
	//tokenList := &user2.TokenResponse{
	//	AccessToken:    token.AccessToken,
	//	RefreshToken:   token.RefreshToken,
	//	TokenType:      "bearer",
	//	AccessTokenExp: token.AccessExp,
	//}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.userRepo.UpdateLoginTime(conn, c, mem.Username)
		if err != nil {
			zap.L().Error("evn_user user_service Login UpdateLoginTime DB_Error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	//err = ls.memberRepo.UpdateLoginTime(c, int64(mem.ID))
	//if err != nil {
	//	zap.L().Error("登陆模块user数据库登陆时间存入出错", zap.Error(err))
	//	return &user2.LoginResponse{}, errs.GrpcError(model.DBError)
	//}
	userInfo := response.UserInfoResponse(mem, token.AccessToken, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	//rsp := &user2.LoginResponse{}
	//err = copier.Copy(&rsp, userInfo)
	//if err != nil {
	//	zap.L().Error("evn_user user_service Login copier.Copy Copy_Error", zap.Error(err))
	//	return &user2.LoginResponse{}, errs.GrpcError(model.CopyError)
	//}
	return &user2.UserInfoResponse{
		Id:        int64(userInfo.ID),
		Username:  userInfo.UserName,
		Photo:     userInfo.Photo,
		Token:     userInfo.Token,
		CreatedAt: userInfo.CreatedAt.String(),
	}, nil
}

func (ls *UserService) Forget(ctx context.Context, req *user2.ForgetRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//可以校验参数
	//校验验证码
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+req.Email)
	if err == redis.Nil {
		return nil, errs.GrpcError(model.CaptchaNoExist)
	}
	if err != nil {
		zap.L().Error("evn_user user_service Forget redis read err", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != req.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	//校验业务逻辑
	exist, err := ls.userRepo.IsExistByName(c, req.Email)
	if err != nil {
		zap.L().Error("evn_user user_service Forget GetUserByEmail DB_Error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if !exist {
		return nil, errs.GrpcError(model.AccountNoExist)
	}
	//执行业务逻辑
	//pwd := encrypts.Md5(req.Password) //加密部分
	//随机生成用户ID
	//使用薄雾算法生成user id
	//mist := common.NewMist()
	//userIdSequence := mist.Generate()

	//bcrypt 密码加密
	pwHashBytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	//转换成字符串
	pwHashStr := string(pwHashBytes)

	mem := &user.User{
		Email:    req.Email,
		Password: pwHashStr,
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.userRepo.UpdateUser(conn, c, mem)
		if err != nil {
			zap.L().Error("evn_user user_service Forget UpdateUser DB_Error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	return &user2.CommonDataResponse{Data: "发送成功"}, nil
}

// TokenVerify token验证
func (ls *UserService) TokenVerify(ctx context.Context, msg *user2.TokenRequest) (*user2.TokenVerifyResponse, error) {
	c := context.Background()
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	//此处为了方便复用，增加一个参数用于接收解析jwt的密钥
	parseToken, err := jwts.ParseToken(token, msg.Secret)
	if err != nil {
		zap.L().Error("evn_user user_service TokenVerify ParseToken error", zap.Error(err))
		return nil, errs.GrpcError(model.NoLoginError)
	}
	//数据库查询 优化点 登陆之后应该把用户信息缓存起来
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	memberById, err := ls.userRepo.FindUserById(c, id)
	if err != nil {
		zap.L().Error("evn_user user_service TokenVerify FindUserById DBError", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	memMessage := &user2.TokenVerifyResponse{
		Id:       int64(memberById.ID),
		Username: memberById.Username,
	}
	if msg.IsEncrypt {
		tmp, _ := encrypts.EncryptInt64(int64(memberById.ID), config.C.AC.AesKey)
		memMessage.Id, _ = strconv.ParseInt(tmp, 10, 64) //加密用户ID
	}
	return memMessage, nil
}

func (ls *UserService) RefreshToken(ctx context.Context, req *user2.RefreshTokenRequest) (*user2.TokenResponse, error) {
	c := context.Background()
	//接收参数
	reqStruct := &user2.TokenRequest{
		Token:     req.RefreshToken,
		Secret:    config.C.JC.RefreshSecret,
		IsEncrypt: false, //不加密 返回的用户ID
	}
	//校验参数
	parseRsp, err := ls.TokenVerify(c, reqStruct)
	if err != nil {
		return nil, err //失败则返回空
	}
	//正确则重新生成Token列表返回
	memId := parseRsp.Id
	//使用jwt生成token
	memIdStr := strconv.FormatInt(memId, 10)
	token := jwts.CreateToken(memIdStr, config.C.JC.AccessExp, config.C.JC.AccessSecret, config.C.JC.RefreshSecret, config.C.JC.RefreshExp)
	tokenList := &user2.TokenResponse{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
		AccessTokenExp: token.AccessExp,
	}
	return tokenList, nil
}

func (ls *UserService) GetSpaceIndividual(ctx context.Context, req *user2.SpaceIndividualRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//var userInfo *user.User
	userInfo, err := ls.userRepo.FindUserById(c, int64(req.ID))
	if err != nil {
		zap.L().Error("evn_user user_service GetSpaceIndividual FindUserById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	isAttention := false
	if req.Uid != 0 {
		//获取是否关注
		isAttention, err = ls.userRepo.IsAttention(c, req.Uid, req.ID)
		if err != nil {
			zap.L().Error("evn_user user_service GetSpaceIndividual IsAttention error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
	}
	//获取关注和粉丝
	attentionNum, err := ls.userRepo.GetAttentionNum(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetSpaceIndividual GetAttentionNum error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	vermicelliNum, err := ls.userRepo.GetVermicelliNum(c, req.ID)

	rsp, err := response.GetSpaceIndividualResponse(userInfo, isAttention, attentionNum, vermicelliNum, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetSpaceIndividual GetSpaceIndividualResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetSpaceIndividual rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetReleaseInformation(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	videoList, err := ls.userRepo.GetVideoListBySpace(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetReleaseInformation GetVideoListBySpace error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	articleList, err := ls.userRepo.GetArticleBySpace(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetReleaseInformation GetArticleBySpace error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp, err := response.GetReleaseInformationResponse(videoList, articleList, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetReleaseInformation GetReleaseInformationResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetReleaseInformation rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetAttentionList(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取用户关注列表
	attentionList, err := ls.userRepo.GetAttentionList(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetAttentionList GetAttentionList error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//获取自己关注的用户
	arr, err := ls.userRepo.GetAttentionListByIdArr(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetAttentionList GetAttentionListByIdArr error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	rsp, err := response.GetAttentionListResponse(attentionList, arr, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetAttentionList GetAttentionListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetAttentionList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetVermicelliList(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取用户粉丝列表
	vermicelliList, err := ls.userRepo.GetVermicelliList(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetVermicelliList GetVermicelliList error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//获取自己关注的用户
	arr, err := ls.userRepo.GetAttentionListByIdArr(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetVermicelliList GetAttentionListByIdArr error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	rsp, err := response.GetVermicelliListResponse(vermicelliList, arr, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetVermicelliList GetVermicelliListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetVermicelliList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetUserInfo(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取用户粉丝列表
	tmpUser, err := ls.userRepo.FindUserById(c, int64(req.ID))
	if err != nil {
		zap.L().Error("evn_user user_service GetUserInfo GetUserByID error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	rsp := response.UserSetInfoResponse(tmpUser)
	if err != nil {
		zap.L().Error("evn_user user_service GetUserInfo UserSetInfoResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetUserInfo rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) SetUserInfo(ctx context.Context, req *user2.UserInfoRequest) (*user2.CommonBoolResponse, error) {
	c := context.Background()
	update := map[string]interface{}{
		"Username":  req.Username,
		"Gender":    0,
		"BirthDate": req.Birth_Date,
		"IsVisible": conversion.BoolTurnInt8(req.Is_Visible),
		"Signature": req.Signature,
	}
	tmpRsp, err := ls.userRepo.UpdatePureZero(c, int64(req.ID), update)
	if err != nil {
		zap.L().Error("evn_user user_service SetUserInfo UpdatePureZero error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if !tmpRsp {
		return nil, errs.GrpcError(model.DBError)
	}
	return &user2.CommonBoolResponse{Data: true}, nil
}

func (ls *UserService) DetermineNameExists(ctx context.Context, req *user2.DetermineNameExistsRequest) (*user2.CommonBoolResponse, error) {
	c := context.Background()
	is, err := ls.userRepo.IsExistByName(c, req.Username)
	if err != nil {
		zap.L().Error("evn_user user_service DetermineNameExists IsExistByName error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if is {
		tmpUser, err := ls.userRepo.FindUserByName(c, req.Username)
		if err != nil {
			zap.L().Error("evn_user user_service DetermineNameExists FindUserByName error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		//判断是否未更改
		if tmpUser.ID == uint(req.ID) {
			return &user2.CommonBoolResponse{Data: false}, nil
		}
		return &user2.CommonBoolResponse{Data: true}, nil
	} else {
		return &user2.CommonBoolResponse{Data: false}, nil
	}

}

func (ls *UserService) UpdateAvatar(ctx context.Context, req *user2.UpdateAvatarRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	photo, _ := json.Marshal(common2.Img{
		Src: req.ImgUrl,
		Tp:  req.TP,
	})
	user := &user.User{PublicModel: common2.PublicModel{ID: uint(req.ID)}, Photo: photo}
	is, err := ls.userRepo.UpdateUserAvatar(c, user)
	if err != nil {
		zap.L().Error("evn_user user_service UpdateAvatar UpdateUserAvatar error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if is {
		url, err := conversion.SwitchIngStorageFun(req.TP, req.ImgUrl, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
		if err != nil {
			zap.L().Error("evn_user user_service UpdateAvatar SwitchIngStorageFun error", zap.Error(err))
			return nil, errs.GrpcError(model.SystemError)
		}
		return &user2.CommonDataResponse{Data: url}, nil
	} else {
		return &user2.CommonDataResponse{Data: "更新失败"}, nil
	}
}

func (ls *UserService) GetLiveData(ctx context.Context, req *user2.CommonIDRequest) (*user2.LiveDataResponse, error) {
	c := context.Background()
	is, err := ls.userRepo.IsExistByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetLiveData IsExistByID error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if is {
		tmpLiveInfo, err := ls.userRepo.FindUserLiveInfo(c, int64(req.ID))
		if err != nil {
			zap.L().Error("evn_user user_service GetLiveData FindUserLiveInfo error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if tmpLiveInfo == nil {
			return &user2.LiveDataResponse{}, nil
		}
		rsp, err := response.GetLiveDataResponse(tmpLiveInfo, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
		if err != nil {
			zap.L().Error("evn_user user_service GetLiveData GetLiveDataResponse error", zap.Error(err))
			return &user2.LiveDataResponse{}, errs.GrpcError(model.SystemError)
		}
		return &user2.LiveDataResponse{
			Title: rsp.Title,
			Img:   rsp.Img,
		}, nil
	}
	return &user2.LiveDataResponse{}, nil
}

func (ls *UserService) SaveLiveData(ctx context.Context, req *user2.SaveLiveDataRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	img, _ := json.Marshal(common2.Img{
		Src: req.Img,
		Tp:  req.TP,
	})
	tmpLiveinfo := &liveInfo.LiveInfo{
		Uid:   uint(req.ID),
		Title: req.Title,
		Img:   datatypes.JSON(img),
	}

	is, err := ls.userRepo.IsExistUserLiveInfo(c, tmpLiveinfo)
	if err != nil {
		zap.L().Error("evn_user user_service SaveLiveData UpdateUserLiveInfo error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if is {
		ok, err := ls.userRepo.UpdateUserLiveInfo(c, tmpLiveinfo)
		if err != nil {
			zap.L().Error("evn_user user_service SaveLiveData UpdateUserLiveInfo error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if ok {
			return &user2.CommonDataResponse{Data: "修改成功"}, nil
		}

	} else {
		ok, err := ls.userRepo.SaveUserLiveInfo(c, tmpLiveinfo)
		if err != nil {
			zap.L().Error("evn_user user_service SaveLiveData SaveUserLiveInfo error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if ok {
			return &user2.CommonDataResponse{Data: "修改成功"}, nil
		}
	}
	return &user2.CommonDataResponse{}, errs.GrpcError(model.SystemError)
}

func (ls *UserService) SendEmailVerificationCodeByChangePassword(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	//1.获取参数
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
	defer cancel()
	tmpUser, err := ls.userRepo.FindUserById(c, int64(req.ID))
	if err != nil {
		zap.L().Error("evn_user user_service SendEmailVerificationCodeByChangePassword FindUserById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	//3.生成验证码(随机四位1000-9999或者六位100000-999999)
	code := util.CreateCaptcha(6) //生成随机六位数字验证码
	fmt.Printf("%v验证码为：%v", tmpUser.Email, code)
	//4.调用短信平台(第三方 放入go func 协程 接口可以快速响应
	//TODO 完善邮件发送服务
	go func() {
		//发送方
		mailTo := []string{tmpUser.Email}
		// 邮件主题 //TODO 测试临时打印验证码
		subject := "验证码"
		// 邮件正文
		body := fmt.Sprintf("您正在找回密码您的验证码为:%s,5分钟有效,请勿转发他人", code)
		//TODO 测试暂时忽略错误
		_ = email1.SendMail(mailTo, subject, body)

		c2, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
		defer cancel()
		err := ls.cache.Put(c2, model.ChangeRedisKey+tmpUser.Email, code, 5*time.Minute)
		if err != nil {
			zap.L().Error("evn_user user_service SendEmailVerificationCodeByChangePassword redis put err", zap.Error(err))

		}
		//zap.L().Debug("将手机号和验证码存入redis成功：REGISTER_" + email + " : " + code + "\n")
	}()
	//注意code一般不发送
	//这里是做了简化处理 由于短信平台目前对于个人不好使用
	return &user2.CommonDataResponse{Data: "发送成功"}, nil
}

func (ls *UserService) ChangePassword(ctx context.Context, req *user2.ChangePasswordRequest) (*user2.CommonDataResponse, error) {
	//可以校验参数
	//校验验证码
	//1.获取参数
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
	defer cancel()
	tmpUser, err := ls.userRepo.FindUserById(c, int64(req.ID))
	if err != nil {
		zap.L().Error("evn_user user_service ChangePassword FindUserById error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	redisCode, err := ls.cache.Get(c, model.ChangeRedisKey+tmpUser.Email)
	if err == redis.Nil {
		return nil, errs.GrpcError(model.CaptchaNoExist)
	}
	if err != nil {
		zap.L().Error("evn_user user_service ChangePassword redis read err", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != req.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}

	//bcrypt 密码加密
	pwHashBytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	//转换成字符串
	pwHashStr := string(pwHashBytes)

	mem := &user.User{
		Email:    tmpUser.Email,
		Password: pwHashStr,
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.userRepo.UpdateUser(conn, c, mem)
		if err != nil {
			zap.L().Error("evn_user user_service ChangePassword UpdateUser DB_Error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})
	return &user2.CommonDataResponse{Data: "修改成功"}, nil
}

func (ls *UserService) Attention(ctx context.Context, req *user2.CommonIDAndUIDRequest) (*user2.CommonDataResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
	defer cancel()
	att, err := ls.userRepo.Attention(c, req.ID, req.UID)
	if err != nil {
		zap.L().Error("evn_user user_service Attention Attention error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if att.Uid == uint(req.UID) {
		return &user2.CommonDataResponse{Data: "操作失败"}, nil
	} else {
		return &user2.CommonDataResponse{Data: "操作成功"}, nil
	}
}

func (ls *UserService) CreateFavorites(ctx context.Context, req *user2.FavoritesRequest) (*user2.CommonDataResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
	defer cancel()
	if req.ID == 0 {
		//插入
		if len(req.Title) == 0 {
			return &user2.CommonDataResponse{Data: "标题为空"}, nil
		}
		//判断是否只有标题
		if req.ID <= 0 && len(req.Tp) == 0 && len(req.Content) == 0 && len(req.Cover) == 0 {
			//单标题创建
			fs := &favorites.Favorites{Uid: uint(req.Uid), Title: req.Title, Max: 1000}
			is, err := ls.userRepo.SaveFavorites(c, fs)
			if !is {
				zap.L().Error("evn_user user_service CreateFavorites SaveFavorites error", zap.Error(err))
				return &user2.CommonDataResponse{Data: "创建失败"}, errs.GrpcError(model.DBError)
			} else {
				return &user2.CommonDataResponse{Data: "创建成功"}, nil
			}
		} else {
			//资料齐全
			cover, _ := json.Marshal(common2.Img{
				Src: req.Cover,
				Tp:  req.Tp,
			})
			fs := &favorites.Favorites{
				Uid:     uint(req.Uid),
				Title:   req.Title,
				Content: req.Content,
				Cover:   cover,
				Max:     1000,
			}
			is, err := ls.userRepo.SaveFavorites(c, fs)
			if !is {
				zap.L().Error("evn_user user_service CreateFavorites SaveFavorites error", zap.Error(err))
				return &user2.CommonDataResponse{Data: "创建失败"}, errs.GrpcError(model.DBError)
			} else {
				return &user2.CommonDataResponse{Data: "创建成功"}, nil
			}
		}
	} else {
		//更新
		fs, err := ls.userRepo.FindFavoritesByID(c, req.ID)
		if err != nil {
			zap.L().Error("evn_user user_service CreateFavorites FindFavoritesByID error", zap.Error(err))
			return &user2.CommonDataResponse{Data: "查询失败"}, errs.GrpcError(model.DBError)
		}
		if fs.Uid != uint(req.Uid) {
			return &user2.CommonDataResponse{Data: "查询非法操作"}, errs.GrpcError(model.DBError)
		}
		cover, _ := json.Marshal(common2.Img{
			Src: req.Cover,
			Tp:  req.Tp,
		})
		fs.Title = req.Title
		fs.Content = req.Content
		fs.Cover = cover
		if is, err := ls.userRepo.UpdateFavorities(c, fs); !is {
			zap.L().Error("evn_user user_service CreateFavorites UpdateFavorities error", zap.Error(err))
			return &user2.CommonDataResponse{Data: "更新失败"}, errs.GrpcError(model.DBError)
		}
		return &user2.CommonDataResponse{Data: "更新成功"}, errs.GrpcError(model.DBError)
	}

}

func (ls *UserService) GetFavoritesList(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取收藏夹列表
	fl, err := ls.userRepo.GetFavoritesList(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoritesList GetFavoritesList error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp, err := response.GetFavoritesListResponse(fl, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoritesList GetFavoritesListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoritesList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) DeleteFavorites(ctx context.Context, req *user2.CommonIDAndUIDRequest) (*user2.CommonDataResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
	defer cancel()
	fs, err := ls.userRepo.FindFavoritesByID(c, req.ID)
	if err != nil {
		zap.L().Error("evn_user user_service DeleteFavorites FindFavoritesByID error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if fs.ID <= 0 {
		return &user2.CommonDataResponse{Data: "收藏夹不存在"}, nil
	}
	if is, err := ls.userRepo.DeleteFavorites(c, fs); is && err == nil {
		if fs.Uid != uint(req.UID) {
			return &user2.CommonDataResponse{Data: "非创建者不可删除"}, nil
		}
	} else {
		zap.L().Error("evn_user user_service DeleteFavorites DeleteFavorites error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "删除失败"}, errs.GrpcError(model.DBError)
	}
	//删除收藏记录
	if is, err := ls.userRepo.DetectCollectByFavoritesID(c, req.ID); is && err == nil {
		return &user2.CommonDataResponse{Data: "删除成功"}, nil
	} else {
		zap.L().Error("evn_user user_service DeleteFavorites DetectCollectByFavoritesID error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "删除失败"}, errs.GrpcError(model.DBError)
	}
}

func (ls *UserService) FavoriteVideo(ctx context.Context, req *user2.FavoriteVideoRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取收藏夹
	for _, v := range req.IDs {
		fl, err := ls.userRepo.FindFavoritesByID(c, v)
		if err != nil {
			zap.L().Error("evn_user user_service FavoriteVideo FindFavoritesByID error", zap.Error(err))
			return nil, errs.GrpcError(model.DBError)
		}
		if fl.Uid != uint(req.UID) {
			return &user2.CommonDataResponse{Data: "非法操作"}, nil
		}
		if len(fl.CollectList) > fl.Max {
			return &user2.CommonDataResponse{Data: "收藏夹已满"}, nil
		}
		cl := &collect.Collect{
			Uid:         uint(req.UID),
			FavoritesID: uint(v),
			VideoID:     uint(req.Video_ID),
		}
		if is, err := ls.userRepo.SaveCollect(c, cl); !is {
			zap.L().Error("evn_user user_service FavoriteVideo SaveCollect error", zap.Error(err))
			return &user2.CommonDataResponse{Data: "收藏失败"}, errs.GrpcError(model.DBError)
		}
	}

	return &user2.CommonDataResponse{Data: "操作成功"}, nil
}

func (ls *UserService) GetFavoritesListByFavoriteVideo(ctx context.Context, req *user2.FavoritesListRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取收藏夹列表
	fl, err := ls.userRepo.GetFavoritesList(c, req.UID)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoritesListByFavoriteVideo GetFavoritesList error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	//查询该视频在那些收藏夹内已收藏
	cl, err := ls.userRepo.FindVideoExistWhere(c, req.Video_ID)
	if cl == nil || err != nil {
		zap.L().Error("evn_user user_service GetFavoritesListByFavoriteVideo FindVideoExistWhere error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	ids := make([]uint, 0)
	for _, v := range *cl {
		ids = append(ids, v.FavoritesID)
	}
	rsp, err := response.GetFavoritesListByFavoriteVideoResponse(fl, ids, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoritesListByFavoriteVideo GetFavoritesListByFavoriteVideoResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoritesListByFavoriteVideo rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetFavoriteVideoList(ctx context.Context, req *user2.FavoriteVideoListRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取收藏夹内视频列表
	cl, err := ls.userRepo.GetVideoInfoByFavoriteID(c, req.Favorite_ID)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoriteVideoList GetVideoInfoByFavoriteID error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	rsp, err := response.GetFavoriteVideoListResponse(cl, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoriteVideoList GetFavoriteVideoListResponse error", zap.Error(err))
		return nil, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetFavoriteVideoList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetRecordList(ctx context.Context, req *user2.GetRecordListRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取历史记录列表
	rl, err := ls.userRepo.GetRecordListByUid(c, req)
	if err != nil {
		zap.L().Error("evn_user user_service GetRecordList GetRecordListByUid error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "查询失败"}, errs.GrpcError(model.DBError)
	}

	rsp, err := response.GetRecordListResponse(rl, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetRecordList GetRecordListResponse error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "响应失败"}, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetRecordList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) ClearRecord(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	if is, err := ls.userRepo.ClearRecord(c, req.ID); !is || err != nil {
		zap.L().Error("evn_user user_service ClearRecord ClearRecord error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "清空失败"}, errs.GrpcError(model.DBError)
	}
	return &user2.CommonDataResponse{Data: "清空完成"}, nil
}

func (ls *UserService) DeleteRecordByID(ctx context.Context, req *user2.CommonIDAndUIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	if is, err := ls.userRepo.DeleteRecordByID(c, req); !is || err != nil {
		zap.L().Error("evn_user user_service DeleteRecordByID DeleteRecordByID error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "删除失败"}, errs.GrpcError(model.DBError)
	}
	return &user2.CommonDataResponse{Data: "删除成功"}, nil
}

func (ls *UserService) GetNoticeList(ctx context.Context, req *user2.GetNoticeListRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取用户通知
	messageType := make([]string, 0)
	switch req.Tp {
	case "comment":
		messageType = append(messageType, notice.VideoComment, notice.ArticleComment)
		break
	case "like":
		messageType = append(messageType, notice.VideoLike, notice.ArticleLike)
	}

	//获取历史记录列表
	nl, err := ls.userRepo.GetNoticeList(c, req, messageType)
	if err != nil {
		zap.L().Error("evn_user user_service GetNoticeList GetNoticeList error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "查询失败"}, errs.GrpcError(model.DBError)
	}

	if is, err := ls.userRepo.ReadAllNoticeList(c, req); !is || err != nil {
		zap.L().Error("evn_user user_service GetNoticeList ReadAllNoticeList error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "已读消息失败"}, errs.GrpcError(model.DBError)
	}
	rsp, err := response.GetNoticeListResponse(nl, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetNoticeList GetNoticeListResponse error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "响应失败"}, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetNoticeList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetChatList(ctx context.Context, req *user2.CommonIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取聊天记录列表
	cl, err := ls.userRepo.GetChatListByIO(c, req)
	if err != nil {
		zap.L().Error("evn_user user_service GetChatList GetChatListByIO error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "查询失败"}, errs.GrpcError(model.DBError)
	}
	ids := make([]uint, 0)
	for _, v := range *cl {
		ids = append(ids, v.Tid)
	}
	msgList := make(map[uint]*chatMsg.MsgList, 0)
	for _, v := range ids {
		ml, err := ls.userRepo.FindMsgList(c, req, v)
		if err != nil {
			break
		}
		msgList[v] = ml
	}
	rsp, err := response.GetChatListResponse(cl, msgList, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetChatList GetChatListResponse error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "响应失败"}, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetChatList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetChatHistoryMsg(ctx context.Context, req *user2.GetChatHistoryMsgRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	//获取聊天记录列表
	cm, err := ls.userRepo.FindHistoryMsg(c, req)
	if err != nil {
		zap.L().Error("evn_user user_service GetChatHistoryMsg FindHistoryMsg error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "查询失败"}, errs.GrpcError(model.DBError)
	}
	rsp, err := response.GetChatHistoryMsgResponse(cm, config.C.Host.LocalHost, config.C.Host.TencentOssHost)
	if err != nil {
		zap.L().Error("evn_user user_service GetChatHistoryMsg GetChatHistoryMsgResponse error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "响应失败"}, errs.GrpcError(model.SystemError)
	}
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		zap.L().Error("evn_user user_service GetChatHistoryMsg rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.CommonDataResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) PersonalLetter(ctx context.Context, req *user2.CommonIDAndUIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()
	cm, err := ls.userRepo.GetLastMessage(c, req)
	if err != nil {
		zap.L().Error("evn_user user_service PersonalLetter GetLastMessage error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "操作失败"}, errs.GrpcError(model.DBError)
	}
	var lastTime time.Time
	if cm.ID > 0 {
		lastTime = cm.CreatedAt
	} else {
		lastTime = time.Now()
	}
	ci := &chatList.ChatsListInfo{
		Uid:         uint(req.UID),
		Tid:         uint(req.ID),
		LastMessage: cm.Message,
		LastAt:      lastTime,
	}
	if is, err := ls.userRepo.AddChat(c, ci); !is || err != nil {
		zap.L().Error("evn_user user_service PersonalLetter AddChat error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "操作失败"}, errs.GrpcError(model.DBError)
	}
	return &user2.CommonDataResponse{Data: "操作成功"}, nil
}

func (ls *UserService) DeleteChatItem(ctx context.Context, req *user2.CommonIDAndUIDRequest) (*user2.CommonDataResponse, error) {
	c := context.Background()

	if is, err := ls.userRepo.DeleteChat(c, req); !is || err != nil {
		zap.L().Error("evn_user user_service DeleteChatItem DeleteChat error", zap.Error(err))
		return &user2.CommonDataResponse{Data: "删除失败"}, errs.GrpcError(model.DBError)
	}
	return &user2.CommonDataResponse{Data: "删除成功"}, nil
}
