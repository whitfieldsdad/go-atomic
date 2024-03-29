package atomic

import (
	"log"
	"os/user"

	"github.com/shirou/gopsutil/v3/host"
)

type User struct {
	Artifact `json:",inline"`
	UserId   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	HomeDir  string `json:"home_dir"`
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

func GetUserId() string {
	u, err := GetCurrentUser()
	if err != nil {
		log.Fatalf("Failed to lookup current user: %s", err)
	}
	return u.Id
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
	return &User{
		Artifact: Artifact{
			Id: calculateUserId(u.Uid),
		},
		UserId:   u.Uid,
		Name:     u.Name,
		Username: u.Username,
		HomeDir:  u.HomeDir,
	}
}

func calculateUserId(uid string) string {
	return NewUUID5(GetHostId(), []byte(uid))
}
