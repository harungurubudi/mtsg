package provider

import "github.com/harungurubudi/mtsg/internal/usecase"

type Container struct {
	Authentication usecase.Authentication
}

func ProvideContainer(
	auth usecase.Authentication,
) Container {
	return Container{
		Authentication: auth,
	}
}
