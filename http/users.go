
package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	fbErrors "github.com/filebrowser/filebrowser/v2/errors"
	gpfs "github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/users"
)

var (
	NonModifiableFieldsForNonAdmin = []string{"Username", "Scope", "LockPassword", "Perm", "Commands", "Rules"}
)

type modifyUserRequest struct {
	modifyRequest
	Data *users.User `json:"data"`
}

func getUserID(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	i, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(i), err
}

func getUser(_ http.ResponseWriter, r *http.Request) (*modifyUserRequest, error) {
	if r.Body == nil {
		return nil, fbErrors.ErrEmptyRequest
	}

	req := &modifyUserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return nil, err
	}

	if req.What != "user" {
		return nil, fbErrors.ErrInvalidDataType
	}

	return req, nil
}

func withSelfOrAdmin(fn handleFunc) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		id, err := getUserID(r)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if d.user.ID != id && !d.user.Perm.Admin {
			return http.StatusForbidden, nil
		}

		d.raw = id
		return fn(w, r, d)
	})
}

var usersGetHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	users, err := d.store.Users.Gets(d.server.Root)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, u := range users {
		u.Password = ""
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return renderJSON(w, r, users)
})

var userGetHandler = withSelfOrAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	u, err := d.store.Users.Get(d.server.Root, d.raw.(uint))
	if errors.Is(err, fbErrors.ErrNotExist) {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	//查询目录限额
	err, quota := gpfs.QueryFilesetQuota(u.Username, d.settings.Gpfs)
	if err != nil {
		fmt.Println("Error:", err)
		return http.StatusInternalServerError, err
	}
	if len(quota) != 0 {
		u.QuotaSoft = int64(quota[0].BlockQuota) / 1024 / 1024
		u.QuotaHard = int64(quota[0].BlockLimit) / 1024 / 1024
	}

	u.Password = ""
	if !d.user.Perm.Admin {
		u.Scope = ""
	}
	return renderJSON(w, r, u)
})

var userDeleteHandler = withSelfOrAdmin(func(_ http.ResponseWriter, _ *http.Request, d *data) (int, error) {
	err := d.store.Users.Delete(d.raw.(uint))
	if err != nil {
		return errToStatus(err), err
	}

	return http.StatusOK, nil
})

var userPostHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req, err := getUser(w, r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if len(req.Which) != 0 {
		return http.StatusBadRequest, nil
	}

	if req.Data.Password == "" {
		return http.StatusBadRequest, fbErrors.ErrEmptyPassword
	}
	// 保存用户的密码
	rawPasswd := req.Data.Password

	req.Data.Password, err = users.HashPwd(req.Data.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// 默认创建的路径为./username
	homedir := "./" + req.Data.Username

	req.Data.Scope = homedir

	// 在GPFS Linux中创建用户
	err = d.settings.CreateUser(req.Data.Username, rawPasswd, d.server.Root, d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// 在GPFS中创建用户的目录
	err = gpfs.CreateFileset(req.Data.Username, d.settings.Gpfs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	//需要暂停一下，等待GPFS创建目录，不然创建不了文件集
	time.Sleep(5 * time.Second)
	err = gpfs.SetUserHomeDir(req.Data.Username, d.settings.Gpfs)
	if err != nil {
		fmt.Println("Error:", err)
		return http.StatusInternalServerError, err
	}
	// 设置默认限额
	quotaSoft := req.Data.QuotaSoft
	quotaHard := req.Data.QuotaHard
	err = gpfs.SetQuota(req.Data.Username, uint16(quotaSoft), uint16(quotaHard), d.settings.Gpfs)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = d.store.Users.Save(req.Data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Location", "/settings/users/"+strconv.FormatUint(uint64(req.Data.ID), 10))
	return http.StatusCreated, nil
})

var userPutHandler = withSelfOrAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req, err := getUser(w, r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if req.Data.ID != d.raw.(uint) {
		return http.StatusBadRequest, nil
	}

	if len(req.Which) == 0 || (len(req.Which) == 1 && req.Which[0] == "all") {
		if !d.user.Perm.Admin {
			return http.StatusForbidden, nil
		}

		if req.Data.Password != "" {
			req.Data.Password, err = users.HashPwd(req.Data.Password)
		} else {
			var suser *users.User
			suser, err = d.store.Users.Get(d.server.Root, d.raw.(uint))
			req.Data.Password = suser.Password
		}

		if err != nil {
			return http.StatusInternalServerError, err
		}

		req.Which = []string{}
	}

	for k, v := range req.Which {
		v = cases.Title(language.English, cases.NoLower).String(v)
		req.Which[k] = v

		if v == "Password" {
			if !d.user.Perm.Admin && d.user.LockPassword {
				return http.StatusForbidden, nil
			}

			req.Data.Password, err = users.HashPwd(req.Data.Password)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}

		for _, f := range NonModifiableFieldsForNonAdmin {
			if !d.user.Perm.Admin && v == f {
				return http.StatusForbidden, nil
			}
		}
	}
	//修改限额
	err = gpfs.SetQuota(req.Data.Username, uint16(req.Data.QuotaSoft), uint16(req.Data.QuotaHard), d.settings.Gpfs)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = d.store.Users.Update(req.Data, req.Which...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
})
