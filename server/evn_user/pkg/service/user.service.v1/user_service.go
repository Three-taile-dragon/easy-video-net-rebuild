package user_service_v1

import (
	"context"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/encrypts"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/jwts"
	"dragonsss.cn/evn_common/model"
	mCommon "dragonsss.cn/evn_common/model/common"
	"dragonsss.cn/evn_common/model/user"
	user2 "dragonsss.cn/evn_grpc/user"
	"dragonsss.cn/evn_user/config"
	"dragonsss.cn/evn_user/internal/dao"
	"dragonsss.cn/evn_user/internal/dao/mysql"
	"dragonsss.cn/evn_user/internal/database"
	"dragonsss.cn/evn_user/internal/database/tran"
	"dragonsss.cn/evn_user/internal/repo"
	response "dragonsss.cn/evn_user/pkg/model"
	"dragonsss.cn/evn_user/util"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

// LoginService grpc 登陆服务实现
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

func (ls *UserService) GetCaptcha(ctx context.Context, req *user2.CaptchaRequest) (*user2.CaptchaResponse, error) {
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
	go func() {
		//time.Sleep(2 * time.Second)
		//zap.L().Info("短信平台调用成功，发送短信")
		//logs.LG.Debug("短信平台调用成功，发送短信 debug")
		//zap.L().Debug("短信平台调用成功，发送短信 debug")
		//zap.L().Error("短信平台调用成功，发送短信 error")
		//redis存储	假设后续缓存可能存在mysql当中,也可以存在mongo当中,也可能存在memcache当中
		//使用接口 达到低耦合高内聚
		//5.存储验证码 redis 当中,过期时间15分钟
		//redis.Set"REGISTER_"+mobile, code)
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second) //编写上下文 最多允许两秒超时
		defer cancel()
		err := ls.cache.Put(c, "REGISTER_"+email, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("evn_user user_service GetCaptcha redis put err", zap.Error(err))

		}
		//zap.L().Debug("将手机号和验证码存入redis成功：REGISTER_" + email + " : " + code + "\n")
	}()
	//注意code一般不发送
	//这里是做了简化处理 由于短信平台目前对于个人不好使用
	return &user2.CaptchaResponse{Data: "发送成功"}, nil
}

func (ls *UserService) Register(ctx context.Context, req *user2.RegisterRequest) (*user2.RegisterResponse, error) {
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
	exist, err := ls.userRepo.GetUserByEmail(c, req.Email)
	if err != nil {
		zap.L().Error("evn_user user_service Register GetUserByEmail DB_Error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	//检验用户名
	exist, err = ls.userRepo.GetUserByName(c, req.Name)
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
	return &user2.RegisterResponse{
		Id:        int64(mem.ID),
		Username:  mem.Username,
		Photo:     "",
		Token:     token.AccessToken,
		CreatedAt: mem.CreatedAt.String(),
	}, nil
}

func (ls *UserService) Login(ctx context.Context, req *user2.LoginRequest) (*user2.LoginResponse, error) {
	c := context.Background()
	//获取传入参数
	//校验参数
	//校验用户名和邮箱
	exist, err := ls.userRepo.GetUserByName(c, req.Username)
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
	return &user2.LoginResponse{
		Id:        int64(userInfo.ID),
		Username:  userInfo.UserName,
		Photo:     userInfo.Photo,
		Token:     userInfo.Token,
		CreatedAt: userInfo.CreatedAt.String(),
	}, nil
}

func (ls *UserService) Forget(ctx context.Context, req *user2.ForgetRequest) (*user2.ForgetResponse, error) {
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
	exist, err := ls.userRepo.GetUserByEmail(c, req.Email)
	if err != nil {
		zap.L().Error("evn_user user_service Forget GetUserByEmail DB_Error", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if !exist {
		return nil, errs.GrpcError(model.AccountNoExist)
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
	return &user2.ForgetResponse{Data: "发送成功"}, nil
}

// TokenVerify token验证
func (ls *UserService) TokenVerify(ctx context.Context, msg *user2.TokenRequest) (*user2.LoginResponse, error) {
	c := context.Background()
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	//此处为了方便复用，增加一个参数用于接收解析jwt的密钥
	parseToken, err := jwts.ParseToken(token, msg.Secret)
	if err != nil {
		zap.L().Error("Token解析失败", zap.Error(err))
		return nil, errs.GrpcError(model.NoLoginError)
	}
	//数据库查询 优化点 登陆之后应该把用户信息缓存起来
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	memberById, err := ls.userRepo.FindUserById(c, id)
	if err != nil {
		zap.L().Error("Token验证模块member数据库查询出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	memMessage := &user2.LoginResponse{}
	err = copier.Copy(&memMessage, memberById)
	if err != nil {
		zap.L().Error("Token验证模块memMessage赋值错误", zap.Error(err))
		return nil, errs.GrpcError(model.CopyError)
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

func (ls *UserService) GetSpaceIndividual(ctx context.Context, req *user2.SpaceIndividualRequest) (*user2.SpaceIndividualResponse, error) {
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
		zap.L().Error("evn_api user_service GetSpaceIndividual rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.SpaceIndividualResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetReleaseInformation(ctx context.Context, req *user2.ReleaseInformationRequest) (*user2.ReleaseInformationResponse, error) {
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
		zap.L().Error("evn_api user_service GetReleaseInformation rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.ReleaseInformationResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetAttentionList(ctx context.Context, req *user2.AttentionListRequest) (*user2.AttentionListResponse, error) {
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
		zap.L().Error("evn_api user_service GetAttentionList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.AttentionListResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}

func (ls *UserService) GetVermicelliList(ctx context.Context, req *user2.VermicelliListRequest) (*user2.VermicelliListResponse, error) {
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
		zap.L().Error("evn_api user_service GetVermicelliList rspJSON error", zap.Error(err))
		return nil, errs.GrpcError(model.JsonError)
	}
	tmp := &user2.VermicelliListResponse{
		Data: string(rspJSON),
	}
	return tmp, nil
}
