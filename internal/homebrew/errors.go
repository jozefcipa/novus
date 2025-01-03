package homebrew

type HomebrewMissingError struct{}

func (e *HomebrewMissingError) Error() string {
	return "Novus requires Homebrew installed"
}
