package service

import (
	"net/http"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

type Tasks []database.Task

func (tasks Tasks) Len() int {
	return len(tasks)
}

func (tasks Tasks) Less(i, j int) bool {
	return tasks[i].Deadline.Before(tasks[j].Deadline)
}

func (tasks Tasks) Swap(i, j int) {
	tasks[i], tasks[j] = tasks[j], tasks[i]
}

func loginFailed(users []database.User, user database.User) bool {
	for _, u := range users {
		tmp := sha256.Sum256([]byte(user.Pwd))
        hash := hex.EncodeToString(tmp[:])
		if user.Name == u.Name && hash == u.Pwd {
			return false
		}
	}
	return true
}

func contains(users []database.User, new database.User) bool {
	for _, user := range users {
		if new.Name == user.Name {
			return true
		}
	}
	return false
}

func stringToTime(str string) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	time, _ := time.ParseInLocation("2006-01-02 15:04:05", str[:10]+" "+str[11:]+":00", jst)
	return time
}

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	userId, err := ctx.Cookie("user_id")
	if err != nil {
		ctx.HTML(http.StatusOK, "not_login.html", gin.H{"Title": "Not log in"})
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	query, qok := ctx.GetQuery("query")
	tmp, iok := ctx.GetQuery("IsDone")
	isDone, _ := strconv.ParseBool(tmp)

	// Get tasks in DB
	var tasks Tasks
	order := "SELECT tasks.* FROM tasks, task_owners " +
		     "WHERE tasks.id=task_owners.task_id " +
			 "AND task_owners.user_id=" + userId + " "
	if qok {
		order = order + "AND tasks.title LIKE " + "'%" + query + "%'"
	} else if iok && isDone {
		order = order + "AND tasks.is_done=b'1'"
	} else if iok {
		order = order + "AND tasks.is_done=b'0'"
	}
	err = db.Select(&tasks, order) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	dlTmp, _ := ctx.GetQuery("DLOrder")
	dlOrder, _ := strconv.ParseBool(dlTmp)
	if dlOrder {
		sort.Sort(tasks)
	}
	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks, "Query": query, "qok": qok, "IsDone": isDone, "iok": iok})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	_, err := ctx.Cookie("user_id")
	if err != nil {
		ctx.HTML(http.StatusOK, "not_login.html", gin.H{"Title": "Not log in"})
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	tmp := task.Deadline.Format("2006-01-02 15:04:05")
	dlString := tmp[:10] + "T" + tmp[11:16]
	// Render task
	ctx.HTML(http.StatusOK, "task.html", gin.H{"Title": "Task", "Task": task, "DLString": dlString})
}

func Create(ctx *gin.Context) {
	_, err := ctx.Cookie("user_id")
	if err != nil {
		ctx.HTML(http.StatusOK, "not_login.html", gin.H{"Title": "Not log in"})
		return
	}

	ctx.HTML(http.StatusOK, "create.html", gin.H{"Title": "Create"})
}

func CreateConfirm(ctx *gin.Context) {
	userId, _ := ctx.Cookie("user_id")

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var task database.Task
	ctx.Bind(&task)
	dlString, _ := ctx.GetPostForm("DLString")
	if dlString == "" {
		data := map[string]interface{}{"title": task.Title}
		_, err = db.NamedExec("INSERT INTO tasks (title, deadline) VALUES (:title, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 1 YEAR))", data)
	} else {
		deadline := stringToTime(dlString)
        data := map[string]interface{}{"title": task.Title, "deadline": deadline}
		_, err = db.NamedExec("INSERT INTO tasks (title, deadline) VALUES (:title, :deadline)", data)
	}
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var taskId int
	err = db.Get(&taskId, "SELECT LAST_INSERT_ID()")
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ownerData := map[string]interface{}{"task_id": taskId, "user_id": userId}
	_, err = db.NamedExec("INSERT INTO task_owners (task_id, user_id) VALUES (:task_id, :user_id)", ownerData)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "create_confirm.html", gin.H{"Title": "Create confirm", "Task": task})
}

func Edit(ctx *gin.Context) {
	_, err := ctx.Cookie("user_id")
	if err != nil {
		ctx.HTML(http.StatusOK, "not_login.html", gin.H{"Title": "Not log in"})
		return
	}

	dlString, _ := ctx.GetQuery("DLString")
	var task database.Task
	ctx.Bind(&task)
	ctx.HTML(http.StatusOK, "edit.html", gin.H{"Title": "Edit", "Task": task, "DLString": dlString})
}

func EditConfirm(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var task database.Task
	ctx.Bind(&task)
	dlString, _ := ctx.GetPostForm("DLString")
	if dlString == "" {
		data := map[string]interface{}{"id": task.ID, "title": task.Title, "is_done": task.IsDone}
		_, err = db.NamedExec("UPDATE tasks SET title=:title, deadline=DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 1 YEAR), is_done=:is_done WHERE id=:id", data)
	} else {
		deadline := stringToTime(dlString)
        data := map[string]interface{}{"id": task.ID, "title": task.Title, "deadline": deadline, "is_done": task.IsDone}
		_, err = db.NamedExec("UPDATE tasks SET title=:title, deadline=:deadline, is_done=:is_done WHERE id=:id", data)
	}
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "edit_confirm.html", gin.H{"Title": "Edit confirm", "Task": task})
}

func Delete(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var task database.Task
	ctx.Bind(&task)
	data := map[string]interface{}{"id": task.ID}
	_, err = db.NamedExec("DELETE FROM tasks WHERE id=:id", data)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = db.NamedExec("DELETE FROM task_owners WHERE task_id=:id", data)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "delete.html", gin.H{"Title": "Delete", "Task": task})
}

func Share(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var task database.Task
	ctx.Bind(&task)
	name, _ := ctx.GetPostForm("Name")
	var userId int
	err = db.Get(&userId, "SELECT id FROM users WHERE name=?", name)
	if err != nil {
		ctx.HTML(http.StatusOK, "share.html", gin.H{"Title": "Share", "Task": task, "Name": name})
		return
	}
	data := map[string]interface{}{"task_id": task.ID, "user_id": userId}
	_, err = db.NamedExec("INSERT INTO task_owners (task_id, user_id) VALUES (:task_id, :user_id)", data)
	if err != nil {
		ctx.HTML(http.StatusOK, "already_have.html", gin.H{"Title": "Share", "Task": task, "Name": name})
		return
	}

	ctx.HTML(http.StatusOK, "share.html", gin.H{"Title": "Share", "Task": task, "Name": name})
}

func Search(ctx *gin.Context) {
	_, err := ctx.Cookie("user_id")
	if err != nil {
		ctx.HTML(http.StatusOK, "not_login.html", gin.H{"Title": "Not log in"})
		return
	}

	ctx.HTML(http.StatusOK, "search.html", gin.H{"Title": "Search"})
}

func LogIn(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title": "Log in"})
}

func LogInConfirm(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var user database.User
	ctx.Bind(&user)

	// Get users in DB
	var users []database.User
	err = db.Select(&users, "SELECT * FROM users") // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if loginFailed(users, user) {
		ctx.HTML(http.StatusOK, "login_failed.html", gin.H{"Title": "Log in failed"})
		return
	}

	var id int
	err = db.Get(&id, "SELECT id FROM users WHERE name=?", user.Name)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.SetCookie("user_id", strconv.Itoa(id), 3600, "/", "localhost", false, true)
	ctx.HTML(http.StatusOK, "login_confirm.html", gin.H{"Title": "Log in confirm", "Name": user.Name})
}

func EditAccount(ctx *gin.Context) {
	id, err := ctx.Cookie("user_id")
	if err != nil {
		ctx.HTML(http.StatusOK, "not_login.html", gin.H{"Title": "Not log in"})
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "edit_account.html", gin.H{"Title": "Edit account", "User": user})
}

func EditAccountConfirm(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var user database.User
	ctx.Bind(&user)
	// Get users in DB
	var users []database.User
	err = db.Select(&users, "SELECT * FROM users WHERE id<>?", user.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	confirm, _ := ctx.GetPostForm("Confirm")
	if contains(users, user) || user.Pwd != confirm {
		ctx.HTML(http.StatusOK, "name_exist.html", gin.H{"Title": "Name exist", "Name": user.Name})
		return
	}

	tmp := sha256.Sum256([]byte(user.Pwd))
	hash := hex.EncodeToString(tmp[:])
	data := map[string]interface{}{"id": user.ID, "name": user.Name, "pwd": hash}
	_, err = db.NamedExec("UPDATE users SET name=:name, pwd=:pwd WHERE id=:id", data)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "edit_account_confirm.html", gin.H{"Title": "Edit account confirm", "Name": user.Name})
}

func SignUp(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", gin.H{"Title": "Sign up"})
}

func SignUpConfirm(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var new database.User
	ctx.Bind(&new)

	// Get users in DB
	var users []database.User
	err = db.Select(&users, "SELECT * FROM users") // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	confirm, _ := ctx.GetPostForm("Confirm")
	if contains(users, new) || new.Pwd != confirm {
		ctx.HTML(http.StatusOK, "signup_failed.html", gin.H{"Title": "Sign up failed"})
		return
	}
	tmp := sha256.Sum256([]byte(new.Pwd))
	hash := hex.EncodeToString(tmp[:])
	data := map[string]interface{}{"name": new.Name, "pwd": hash}
	_, err = db.NamedExec("INSERT INTO users (name, pwd) VALUES (:name, :pwd)", data)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var userId int
	err = db.Get(&userId, "SELECT LAST_INSERT_ID()")
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.SetCookie("user_id", strconv.Itoa(userId), 3600, "/", "localhost", false, true)
	ctx.HTML(http.StatusOK, "signup_confirm.html", gin.H{"Title": "Sign up confirm", "Name": new.Name})
}

func DeleteAccount(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "delete_account.html", gin.H{"Title": "Delete account"})
}

func DeleteAccountConfirm(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var user database.User
	ctx.Bind(&user)

	// Get tasks in DB
	var users []database.User
	err = db.Select(&users, "SELECT * FROM users") // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if loginFailed(users, user) {
		ctx.HTML(http.StatusOK, "account_not_exist.html", gin.H{"Title": "Account not exist", "Name": user.Name})
		return
	}
	data := map[string]interface{}{"name": user.Name}
	_, err = db.NamedExec("DELETE FROM users WHERE name=:name", data)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.SetCookie("user_id", "", -1, "/", "localhost", false, true)
	ctx.HTML(http.StatusOK, "delete_account_confirm.html", gin.H{"Title": "Delete account confirm", "Name": user.Name})
}