package model

type Department struct {
	ID uint64 `json:"id"`
	Description string `json:"description"`
	Disabled bool `json:"disabled"`
}

func CreateDepartment(department *Department) error {
	statement := "insert into public.departments(description, disabled)VALUES($1, $2);"
	_, err := db.Exec(statement, department.Description, department.Disabled)
	return err
}

