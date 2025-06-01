package replyme

func (m *model) emitTUI(t TUIRequest) {
	m.tuiChan <- t
}

type tuiType uint16

const (
	tuiTypeSelectOne tuiType = iota
	tuiTypeSelectSeveral
	tuiTypeInputText
	tuiTypeInputInt
	tuiTypeInputFile
	tuiTypeConfirm
)

type TUISelectItem struct {
	ID   string
	Name string
	Desc string
}

func (i TUISelectItem) Title() string {
	return i.Name
}

func (i TUISelectItem) Description() string {
	return i.Desc
}

func (i TUISelectItem) FilterValue() string {
	return i.Name
}

type TUISelectOneParams struct {
	Name        string
	Description string
	Items       []TUISelectItem
}

type TUIInputTextParams struct {
	Name        string
	Description string
	Placeholder string
	IsPassword  bool
	Validate    func(s string) bool
	MaxLength   int
}

type TUIInputIntParams struct {
	Name        string
	Description string
	MinValue    int
	MaxValue    int
	Validate    func(s string) bool
}

type TUIInputFileParams struct {
	Name        string
	Description string
	Extensions  []string
	MaxFileSize int
	DoNotOutput bool
}

type TUIConfirmParams struct {
	Name        string
	Description string
}

type TUISelectOneResult struct {
	SelectedID   string
	SelectedItem TUISelectItem
}

type TUIInputFileResult struct {
	Path string
	File []byte
}

type TUIRequest struct {
	ID       string
	Type     tuiType
	Payload  interface{}
	Response chan TUIResponse
}

type TUIResponse struct {
	Value interface{}
	Err   error
}
