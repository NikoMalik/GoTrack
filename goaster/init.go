package goaster

var Goaster *Toaster

func Init() *Toaster {
	Goaster = NewToaster(

		WithBorder(false),
		WithPosition(TopCenter),
		WithVariant(AccentDark),
		WithAutoDismiss(false),
	)
	return Goaster
}
