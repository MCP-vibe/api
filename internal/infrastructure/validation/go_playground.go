package validation

import (
	"api/internal/adapters/validator"
	"errors"
	"regexp"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	go_playground "github.com/go-playground/validator/v10"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
)

type goPlayground struct {
	validator *go_playground.Validate
	translate ut.Translator
	err       error
	msg       []string
}

func NewGoPlayground() (validator.Validator, error) {
	var (
		language         = ru.New()
		uni              = ut.New(language, language)
		translate, found = uni.GetTranslator("ru")
	)

	if !found {
		return nil, errors.New("translator not found")
	}

	v := go_playground.New()
	if err := ru_translations.RegisterDefaultTranslations(v, translate); err != nil {
		return nil, errors.New("translator not found")
	}

	// Зарегистрируйте кастомный валидатор для типа bool
	v.RegisterValidation("bool_required", func(fl go_playground.FieldLevel) bool {
		return true // Простая проверка, так как bool всегда имеет значение
	})

	v.RegisterTranslation("bool_required", translate, func(ut ut.Translator) error {
		return ut.Add("bool_required", "{0} обязательное поле", true)
	}, func(ut ut.Translator, fe go_playground.FieldError) string {
		t, _ := ut.T("bool_required", fe.Field())
		return t
	})

	v.RegisterValidation("phone", validatePhoneNumber)
	v.RegisterTranslation("phone", translate, func(ut ut.Translator) error {
		return ut.Add("phone", "phone должен быть в формате 79xxxxxxxxx", true)
	}, func(ut ut.Translator, fe go_playground.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})

	v.RegisterValidation("password", validatePassword)
	v.RegisterTranslation("password", translate, func(ut ut.Translator) error {
		return ut.Add("password", "password должен быть минимум 8 символов и состоять из букв и цифр", true)
	}, func(ut ut.Translator, fe go_playground.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})

	return &goPlayground{validator: v, translate: translate}, nil
}

func (g *goPlayground) Validate(i interface{}) error {
	if len(g.msg) > 0 {
		g.msg = nil
	}
	g.err = g.validator.Struct(i)
	if g.err != nil {
		return g.err
	}

	return nil
}

func (g *goPlayground) Messages() []string {
	if g.err != nil {
		for _, err := range g.err.(go_playground.ValidationErrors) {
			g.msg = append(g.msg, err.Translate(g.translate))
		}
	}

	return g.msg
}

func validatePhoneNumber(fl go_playground.FieldLevel) bool {
	phone := fl.Field().String()
	re := regexp.MustCompile(`^[7][9]\d{9}$`)
	return re.MatchString(phone)
}

func validatePassword(fl go_playground.FieldLevel) bool {
	password := fl.Field().String()

	// Проверка на минимальную длину
	if len(password) < 8 {
		return false
	}

	// Проверка на наличие латинских букв
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)

	// Проверка на наличие цифр
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	// Проверка на отсутствие русских букв
	hasRussianLetter := regexp.MustCompile(`[А-Яа-я]`).MatchString(password)

	// Пароль должен содержать хотя бы одну букву и одну цифру, и не содержать русские буквы
	return hasLetter && hasDigit && !hasRussianLetter
}
