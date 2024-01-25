package atomic

import (
	"os/user"

	"github.com/charmbracelet/log"
	"github.com/shirou/gopsutil/v3/host"
)

type User struct {
	Id             string   `json:"id"`
	UserId         string   `json:"user_id"`
	Name           string   `json:"name"`
	Username       string   `json:"username"`
	PrimaryGroupId string   `json:"primary_group_id"`
	GroupIds       []string `json:"group_ids"`
	HomeDir        string   `json:"home_dir"`
}

func ListUsers() ([]User, error) {
	users, err := host.Users()
	if err != nil {
		return nil, err
	}
	var result []User
	for _, u := range users {
		username := u.User
		user, err := GetUserByUsername(username)
		if err != nil {
			log.Warnf("Failed to lookup user: %s - %s", username, err)
			continue
		}
		result = append(result, *user)
	}
	return result, nil
}

func GetCurrentUser() (*User, error) {
	o, err := user.Current()
	if err != nil {
		return nil, err
	}
	u := getUser(*o)
	return u, nil
}

func GetUserByUsername(username string) (*User, error) {
	o, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}
	u := getUser(*o)
	return u, nil
}

func getUser(u user.User) *User {
	gids, _ := u.GroupIds()
	return &User{
		Id:             calculateUserId(u.Uid),
		UserId:         u.Uid,
		Name:           u.Name,
		Username:       u.Username,
		HomeDir:        u.HomeDir,
		PrimaryGroupId: u.Gid,
		GroupIds:       gids,
	}
}

func calculateUserId(uid string) string {
	return NewUUID5(GetHostId(), []byte(uid))
}
