package brew

type BrewMissingError struct{}

func (e *BrewMissingError) Error() string {
	return "Novus requires Homebrew installed"
}
