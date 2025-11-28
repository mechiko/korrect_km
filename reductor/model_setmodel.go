package reductor

import (
	"fmt"
	"korrectkm/domain"

	"github.com/mechiko/utility"
)

func Model[T domain.Modeler](src domain.Model) (T, error) {
	var resultNil T
	if instance == nil {
		return resultNil, fmt.Errorf("reductor not init")
	}
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	if pageModel, ok := instance.models[src]; ok {
		if !utility.IsPointer(pageModel) {
			return resultNil, fmt.Errorf("reductor internal error model is pointer")
		}
		if resultOK, ok := pageModel.(T); ok {
			returnModel, err := CloneDeep(resultOK)
			if err != nil {
				return resultNil, fmt.Errorf("reductor clone model %w", err)
			}
			return returnModel, nil
		}
		return resultNil, fmt.Errorf("model wrong type %T", pageModel)
	}
	return resultNil, fmt.Errorf("reductor запрошенной модели нет")
}

// не имеет преимущества видного на данный момент
// то что вызов короче без промежуточного объекта пакета может и хорошо
func SetModel[T domain.Modeler](src T, send bool) error {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	if !utility.IsPointer(src) {
		return fmt.Errorf("reductor: model must be a pointer")
	}
	page := src.Model()
	if !domain.IsValidModel(page.String()) {
		return fmt.Errorf("reductor: model type is invalide")
	}

	storeModel, err := CloneDeep(src)
	if err != nil {
		return fmt.Errorf("reductor: clone error %w", err)
	}
	if !utility.IsPointer(storeModel) {
		return fmt.Errorf("reductor: model clone must be a pointer")
	}
	if instance.models == nil {
		instance.models = make(ModelList)
	}
	instance.models[page] = storeModel
	if !send {
		return nil
	}
	// select-based non-blocking send
	if instance.outStateChan != nil {
		select {
		case instance.outStateChan <- page:
		default:
			// channel full—drop this update
		}
	}
	return nil
}
