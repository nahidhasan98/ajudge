package rank

type rankModel struct {
	OJ          string `bson:"OJ" json:"OJ,omitempty"`
	Username    string `bson:"username" json:"username,omitempty"`
	FullName    string `bson:"fullName" json:"fullName,omitempty"`
	TotalSolved int    `bson:"totalSolved" json:"totalSolved"`
}
