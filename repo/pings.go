package repo

func (r *Repository) ConfigDBPing() (result bool) {
	defer func() {
		if r := recover(); r != nil {
			result = false
		}
	}()
	if err := r.ConfigDB().Ping(); err != nil {
		return false
	}
	return true
}

func (r *Repository) A3DBPing() (result bool) {
	defer func() {
		if r := recover(); r != nil {
			result = false
		}
	}()
	if err := r.A3DB().Ping(); err != nil {
		return false
	}
	return true
}

func (r *Repository) ZnakDBPing() (result bool) {
	defer func() {
		if r := recover(); r != nil {
			result = false
		}
	}()
	if err := r.ZnakDB().Ping(); err != nil {
		return false
	}
	return true
}

func (r *Repository) SelfDBPing() (result bool) {
	defer func() {
		if r := recover(); r != nil {
			result = false
		}
	}()
	if err := r.Self().Ping(); err != nil {
		return false
	}
	return true
}
