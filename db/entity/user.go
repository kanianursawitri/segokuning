package entity

type User struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Password        string `json:"password,omitempty"`
	CredentialType  string `json:"credentialType"`
	CredentialValue string `json:"credentialValue"`
}
