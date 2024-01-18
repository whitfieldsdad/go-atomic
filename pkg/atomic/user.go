package atomic

import "os/user"

type User struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Username       string   `json:"username"`
	PrimaryGroupId string   `json:"primary_group_id"`
	GroupIds       []string `json:"group_ids"`
	HomeDir        string   `json:"home_dir"`
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

func getUser(o user.User) *User {
	gids, _ := o.GroupIds()
	return &User{
		Id:             o.Uid,
		Name:           o.Name,
		Username:       o.Username,
		HomeDir:        o.HomeDir,
		PrimaryGroupId: o.Gid,
		GroupIds:       gids,
	}
}
