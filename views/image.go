package views

func NewCreate(ID uint, putUrl string) (res interface{}, err error) {
	v := map[string]interface{}{
		"id":    ID,
		"putUrl":  putUrl,
	}
	return v, nil
}

