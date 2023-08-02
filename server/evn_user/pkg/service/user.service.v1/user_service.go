package user_service_v1

import (
	"context"
	common "dragonsss.cn/evn_common"
	"dragonsss.cn/evn_common/encrypts"
	"dragonsss.cn/evn_common/errs"
	"dragonsss.cn/evn_common/jwts"
	"dragonsss.cn/evn_grpc/user/login"
	"dragonsss.cn/evn_user/config"
	"dragonsss.cn/evn_user/internal/dao"
	"dragonsss.cn/evn_user/internal/dao/mysql"
	"dragonsss.cn/evn_user/internal/data/user"
	"dragonsss.cn/evn_user/internal/database"
	"dragonsss.cn/evn_user/internal/database/tran"
	"dragonsss.cn/evn_user/internal/repo"
	"dragonsss.cn/evn_user/pkg/model"
	"dragonsss.cn/evn_user/util"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// LoginService grpc 登陆服务实现
type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache       repo.Cache
	userRepo    repo.UserRepo
	transaction tran.Transaction
	memberRepo  repo.MemberRepo
}

func New() *LoginService {
	return &LoginService{
		cache:       dao.Rc,
		userRepo:    mysql.NewUserDao(),
		transaction: dao.NewTransaction(),
		memberRepo:  mysql.NewMemberDao(),
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, req *login.CaptchaRequest) (*login.CaptchaResponse, error) {
	//1.获取参数
	mobile := req.UserMobile
	//2.校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile) //使用自定义错误码进行处理
	}
	//3.生成验证码(随机四位1000-9999或者六位100000-999999)
	code := util.CreateCaptcha(6) //生成随机六位数字验证码
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
		err := ls.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("验证码存入redis出错,cause by : " + err.Error() + "\n")

		}
		zap.L().Debug("将手机号和验证码存入redis成功：REGISTER_" + mobile + " : " + code + "\n")
	}()
	//注意code一般不发送
	//这里是做了简化处理 由于短信平台目前对于个人不好使用
	return &login.CaptchaResponse{Code: code}, nil
}

func (ls *LoginService) Register(ctx context.Context, req *login.RegisterRequest) (*login.RegisterResponse, error) {
	c := context.Background()
	//可以校验参数
	//校验验证码
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+req.Mobile)
	if err == redis.Nil {
		return nil, errs.GrpcError(model.CaptchaNoExist)
	}
	if err != nil {
		zap.L().Error("Register 中 redis 读取错误", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != req.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	//校验业务逻辑
	exist, err := ls.userRepo.GetUserByEmail(c, req.Email)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist)
	}
	//检验用户名
	exist, err = ls.userRepo.GetUserByAccount(c, req.Name)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist)
	}
	//检验手机号
	exist, err = ls.userRepo.GetUserByMobile(c, req.Mobile)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist)
	}
	//执行业务逻辑
	pwd := encrypts.Md5(req.Password) //加密部分
	//随机生成用户ID
	//使用薄雾算法生成user id
	mist := common.NewMist()
	userIdSequence := mist.Generate()
	mem := &user.User{
		Id:            userIdSequence,
		Account:       req.Name,
		Password:      pwd,
		Name:          req.Name,
		Mobile:        req.Mobile,
		Email:         req.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: 0,
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.userRepo.SaveUser(conn, c, mem)
		if err != nil {
			zap.L().Error("注册模块user数据库存入出错", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		////存入组织
		//org := &data.Organization{
		//	Name:       mem.Name + "个人项目",
		//	MemberId:   mem.Id,
		//	CreateTime: time.Now().UnixMilli(),
		//	Personal:   model.Personal,
		//	Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		//}
		//err = ls.organizationRepo.SaveOrganization(conn, c, org)
		//if err != nil {
		//	zap.L().Error("注册模块organization数据库存入失败", zap.Error(err))
		//	return errs.GrpcError(model.DBError)
		//}
		return nil
	})
	//var conn database.DbConn
	//err = ls.userRepo.SaveMember(conn, c, mem)
	////err = ls.memberRepo.SaveMember(c, mem)
	//if err != nil {
	//	zap.L().Error("注册模块user数据库存入出错", zap.Error(err))
	//	return &login.RegisterResponse{}, errs.GrpcError(model.DBError)
	//}
	return &login.RegisterResponse{}, nil
}

func (ls *LoginService) Login(ctx context.Context, req *login.LoginRequest) (*login.LoginResponse, error) {
	c := context.Background()
	//获取传入参数
	//校验参数
	//校验用户名和邮箱
	exist, err := ls.userRepo.GetUserByAccountAndEmail(c, req.Account)
	if err != nil {
		zap.L().Error("数据库出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if !exist {
		return nil, errs.GrpcError(model.AccountNoExist)
	}
	//查询账号密码是否正确
	pwd := encrypts.Md5(req.Password)
	mem, err := ls.userRepo.FindUser(c, req.Account, pwd)
	if err != nil {
		zap.L().Error("登陆模块member数据库查询出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if mem == nil {
		return nil, errs.GrpcError(model.AccountAndPwdError)
	}
	memMessage := &login.MemberMessage{}
	err = copier.Copy(memMessage, mem)
	memMessage.Code, _ = encrypts.EncryptInt64(mem.Id, config.C.AC.AesKey) //加密用户ID
	if err != nil {
		zap.L().Error("登陆模块mem赋值错误", zap.Error(err))
		return nil, errs.GrpcError(model.CopyError)
	}
	//使用jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	token := jwts.CreateToken(memIdStr, config.C.JC.AccessExp, config.C.JC.AccessSecret, config.C.JC.RefreshSecret, config.C.JC.RefreshExp)
	tokenList := &login.TokenResponse{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
		AccessTokenExp: token.AccessExp,
	}
	//将存入部分使用事务包裹 使得可以回滚数据库操作
	//err = ls.transaction.Action(func(conn database.DbConn) error {
	//	err = ls.userRepo.UpdateLoginTime(conn, c, mem.Id)
	//	if err != nil {
	//		zap.L().Error("登陆模块user数据库登陆时间存入出错", zap.Error(err))
	//		return errs.GrpcError(model.DBError)
	//	}
	//	return nil
	//})
	//TODO 这里 未使用事务进行包裹 有问题待处理
	err = ls.memberRepo.UpdateLoginTime(c, mem.Id)
	if err != nil {
		zap.L().Error("登陆模块user数据库登陆时间存入出错", zap.Error(err))
		return &login.LoginResponse{}, errs.GrpcError(model.DBError)
	}
	return &login.LoginResponse{
		Member:    memMessage,
		TokenList: tokenList,
	}, nil
}

// TokenVerify token验证
func (ls *LoginService) TokenVerify(ctx context.Context, msg *login.TokenRequest) (*login.LoginResponse, error) {
	c := context.Background()
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "")
	}
	//此处为了方便复用，增加一个参数用于接收解析jwt的密钥
	parseToken, err := jwts.ParseToken(token, msg.Secret)
	if err != nil {
		zap.L().Error("Token解析失败", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	//数据库查询 优化点 登陆之后应该把用户信息缓存起来
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	memberById, err := ls.userRepo.FindUserById(c, id)
	if err != nil {
		zap.L().Error("Token验证模块member数据库查询出错", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	memMessage := &login.MemberMessage{}
	err = copier.Copy(&memMessage, memberById)
	if err != nil {
		zap.L().Error("Token验证模块memMessage赋值错误", zap.Error(err))
		return nil, errs.GrpcError(model.CopyError)
	}
	if msg.IsEncrypt {
		memMessage.Code, _ = encrypts.EncryptInt64(memberById.Id, config.C.AC.AesKey) //加密用户ID
	}
	return &login.LoginResponse{Member: memMessage}, nil
}

func (ls *LoginService) RefreshToken(ctx context.Context, req *login.RefreshTokenRequest) (*login.TokenResponse, error) {
	c := context.Background()
	//接收参数
	reqStruct := &login.TokenRequest{
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
	memId := parseRsp.Member.Id
	//使用jwt生成token
	memIdStr := strconv.FormatInt(memId, 10)
	token := jwts.CreateToken(memIdStr, config.C.JC.AccessExp, config.C.JC.AccessSecret, config.C.JC.RefreshSecret, config.C.JC.RefreshExp)
	tokenList := &login.TokenResponse{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenType:      "bearer",
		AccessTokenExp: token.AccessExp,
	}
	return tokenList, nil
}
