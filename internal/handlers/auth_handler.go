package handlers

//func SignUp(w http.ResponseWriter, r *http.Request) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	var createUserDto *dto.CreateUserDto
//	err := json.NewDecoder(r.Body).Decode(&createUserDto)
//
//	if err != nil {
//		RespondWithError(w, http.StatusBadRequest,
//			"Invalid request body:  "+err.Error())
//		return
//	}
//
//	defer r.Body.Close()
//
//}
