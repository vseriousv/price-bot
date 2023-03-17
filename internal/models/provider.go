package models

type Provider struct {
	Name, ApiUrl string
}

func (p *Provider) GetName() string {
	return p.Name
}
