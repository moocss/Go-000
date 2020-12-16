package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/moocss/example/internal/service"
)

func user(srv service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(path.Base(r.URL.Path), 10, 64)
		// log.Printf("请求参数:%d\n", id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": http.StatusNotFound,
				"msg": "无效的接口地址",
			})
			return
		}

		// 接口层处理从 service 和 dao 层 透传过来的error
		user, err := srv.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("srv.FindByID err: %+v\n", err)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": http.StatusBadRequest,
				"msg": "此用户不存在",
				"data":  user,
			})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": http.StatusOK,
			"msg": "查询用户成功",
			"data":  user,
		})
		return
	}
}

func NewUserHandler(r *http.ServeMux, srv service.UserService)  {
	r.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "okay")
		},
	))

	// GET /user/1
	r.HandleFunc("/user/", user(srv))
}