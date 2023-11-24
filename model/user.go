package model

type User struct {
	ID uint64 `json:"id"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"password"`
	DepartmentId uint64 `json:"departmentid"`
	Disabled bool `json:"disabled"`
}

func CreateUser(user *User)error{
	statement := "INSERT INTO public.users(name, surname, username, password, departmentid, disabled)VALUES($1, $2, $3, $4, $5, $6);"
	_, err := db.Exec(statement, user.Name, user.Surname, user.Username, user.Password, user.DepartmentId, user.Disabled)
	return err
}

func GetUser(id string)(User, error){
	var user User
	statement := `SELECT * FROM public.users WHERE id = $1;`

	rows, err := db.Query(statement, id)
	if err != nil {
		return User{}, err
	}

	for rows.Next(){
		err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Username, &user.Password, &user.DepartmentId, &user.Disabled)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}