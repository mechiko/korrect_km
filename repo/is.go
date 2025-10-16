package repo

func (r *Repository) IsA3() bool {
	return r.dbs.A3().Exists
}

func (r *Repository) IsConfig() bool {
	return r.dbs.ConfigInfo().Exists
}

func (r *Repository) IsZnak() bool {
	return r.dbs.Znak().Exists
}
