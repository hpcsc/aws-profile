package handlers

type GlobalArguments struct {
	CredentialsFilePath *string
	ConfigFilePath      *string
}

type Handler interface {
	Handle(globalArguments GlobalArguments) (bool, string)
}
