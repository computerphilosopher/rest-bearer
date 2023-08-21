package gittypes

type Repository struct {
	Owner string
	Name  string
}

type Commit struct {
	Repository Repository
	Id         string
}
