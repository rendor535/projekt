package models

type Profile struct {
	UserID        int     `json:"user_id"`
	Age           int     `json:"age"`
	Height        int     `json:"height"`
	Weight        float64 `json:"weight"`
	Gender        string  `json:"gender"`
	ActivityLevel float64 `json:"activity_level"`
	Goal          string  `json:"goal"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
}

func (p *Profile) IsEmpty() bool {
	return p.Age == 0 && p.Height == 0 && p.Weight == 0 &&
		p.Gender == "" && p.ActivityLevel == 0 && p.Goal == ""
}

func (p *Profile) IsValid() bool {
	return p.Age > 0 && p.Height > 0 && p.Weight > 0 &&
		p.Gender != "" && p.ActivityLevel > 0 && p.Goal != ""
}
