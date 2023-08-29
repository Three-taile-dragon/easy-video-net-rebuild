package midd

import (
	"context"
	"dragonsss.cn/evn_common/jwts"
	"dragonsss.com/evn_ws/config"
	"dragonsss.com/evn_ws/internal/dao/mysql"
	"dragonsss.com/evn_ws/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func VerificationTokenAsSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		//升级ws 以便返回消息
		conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			http.NotFound(c.Writer, c.Request)
			c.Abort()
			return
		}
		token := c.Query("token")
		claim, err := jwts.ParseToken(token, config.C.JC.AccessSecret)
		if err != nil {
			response.NotLoginWs(conn, "Token 验证失败")
			_ = conn.Close()
			c.Abort()
			return
		}
		wsRepo := mysql.NewWsDao()
		ctx := context.Background()
		id, _ := strconv.ParseInt(claim, 10, 64)
		user, err := wsRepo.FindUserById(ctx, id)
		if user.ID < 0 || err != nil {
			//没有改用户的情况下
			response.NotLoginWs(conn, "用户异常")
			_ = conn.Close()
			c.Abort()
			return
		}
		c.Set("uid", user.ID)
		c.Set("conn", conn)
		c.Set("currentUserName", user.Username)
		c.Next()
	}
}
