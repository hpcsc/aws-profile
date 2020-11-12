package awsconfig

type Profile struct {
	ProfileName        string
	DisplayProfileName string
	RoleArn            string
	MFASerialNumber    string
	Region             string
	SourceProfile      string
}
