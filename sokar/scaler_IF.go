package sokar

type Scaler interface {
	ScaleBy(amount int) error
}
