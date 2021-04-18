package errors

type ErrorCode int

const (
	ErrorNoAddress ErrorCode = iota + 1
	ErrorNoIntent
	ErrorNoCarsNearby
	ErrorNoCarsheringInUserRegion
	ErrorNoCarsInUserRegion
)

type AliceError struct {
	Code ErrorCode
}

func NewError(code ErrorCode) *AliceError {
	e := new(AliceError)
	e.Code = code
	return e
}

func (e *AliceError) Error() string {
	descr, exists := ErrorMap[e.Code]
	if !exists {
		return "Произошла какая-то ошибка"
	}
	return descr
}

var ErrorMap = map[ErrorCode]string{ //TODO
	ErrorNoAddress:    "Мне не удалось определить адрес",
	ErrorNoIntent:     "Мне не удалось распознать команду",
	ErrorNoCarsNearby: "Мне не удалось найти машины поблизости",
	ErrorNoCarsheringInUserRegion: "К сожалению, в вашем регионе недоступен каршеринг",
	ErrorNoCarsInUserRegion: "К сожалению, в вашем регионе нет доступных машин",
}
