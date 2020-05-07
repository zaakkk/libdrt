package coreMail

type MailStruct struct {
	From     string
	Username string
	Password string
	To       string
	Sub      string
	Msg      []byte
}

func NewMailStruct(from string, username string, password string, to string, sub string, msg []byte) *MailStruct {
	m := new(MailStruct)
	m.From = from
	m.Username = username
	m.Password = password
	m.To = to
	m.Sub = sub
	m.Msg = msg

	return m
}

func (m *MailStruct) SetFrom(from string) *MailStruct {
	m.From = from
	return m
}

func (m *MailStruct) SetUserName(username string) *MailStruct {
	m.Username = username
	return m
}

func (m *MailStruct) SetPassword(password string) *MailStruct {
	m.Password = password
	return m
}

func (m *MailStruct) SetTo(to string) *MailStruct {
	m.To = to
	return m
}

func (m *MailStruct) SetSub(sub string) *MailStruct {
	m.Sub = sub
	return m
}

func (m *MailStruct) SetMsg(msg []byte) *MailStruct {
	m.Msg = msg
	return m
}

func (m *MailStruct) GetFrom() string {
	return m.From
}

func (m *MailStruct) GetUserName() string {
	return m.Username
}

func (m *MailStruct) GetPassword() string {
	return m.Password
}

func (m *MailStruct) GetTo() string {
	return m.To
}

func (m *MailStruct) GetSub() string {
	return m.Sub
}

func (m *MailStruct) GetMsg() []byte {
	return m.Msg
}
