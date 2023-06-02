package system

import "time"

type Metadata struct {
	Author            string    `prompt:"What is your name?" yaml:"author"`
	Email             string    `prompt:"What is your email?" yaml:"email"`
	CompanyUniversity string    `prompt:"What is your company or university?" yaml:"company_university"`
	CreatedAt         time.Time `yaml:"created_at"`
	UpdatedAt         time.Time `yaml:"updated_at"`
}

func (m *Metadata) Update() error {
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Metadata) Create() error {
	m.CreatedAt = time.Now()
	return nil
}
