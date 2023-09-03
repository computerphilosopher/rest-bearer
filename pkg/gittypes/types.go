package gittypes

import "fmt"

type Repository struct {
	Owner string
	Name  string
}

func (repo Repository) ToString() string {
	return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
}

type Commit struct {
	Repository Repository
	Id         string
}

type Remote struct {
	BaseUrl    string
	Repository Repository
}

func (remote Remote) ToString() string {
	return fmt.Sprintf("%s/%s", remote.BaseUrl, remote.Repository.ToString())
}
