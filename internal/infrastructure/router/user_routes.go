package router

import (
	"api/internal/adapters/api/action"
	"api/internal/adapters/repo"
	"api/internal/usecase"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new user
// @Description Создать пользователя
// @Tags user
// @Accept json
// @Produce json
// @Param input body usecase.CreateUserInput true "Create user input"
// @Success 200 {object} nil
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func (s GinEngine) buildCreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateUserUseCase(
				repo.NewUserDB(s.db),
				s.ctxTimeout,
			)
			act = action.NewCreateUserAction(uc, s.log, s.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}
