package reductor

// вернет модель из мап или nil если запрошенной модели нет
func (rdc *Reductor) Model(page ModelType) interface{} {
	rdc.mutex.Lock()
	defer rdc.mutex.Unlock()

	if pageModel, ok := rdc.models[page]; ok {
		return pageModel
	}
	return nil
}

// записываем модель по типу енум моделей
func (rdc *Reductor) SetModel(page ModelType, model interface{}) {
	rdc.mutex.Lock()
	defer rdc.mutex.Unlock()
	if rdc.models == nil {
		rdc.models = make(ModelList)
	}
	rdc.models[page] = model
}
