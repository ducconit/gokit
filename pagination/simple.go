package pagination

type Simple struct {
	Page  int `query:"page" form:"page" json:"page" default:"1"`
	Limit int `query:"limit" form:"limit" json:"limit" default:"20" max:"100"`
}

func (p *Simple) Normalize() {
	Normalize(p)
}

func (p Simple) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}
