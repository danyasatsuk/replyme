package replyme

func (m *model) emitTUI(t TUIRequest) {
	m.tuiChan <- t
}

type tuiType uint16

const (
	tuiType_SelectOne tuiType = iota
	tuiType_SelectSeveral
	tuiType_InputText
	tuiType_InputInt
	tuiType_InputFile
	tuiType_Confirm
)

type tuiSelectItem struct {
	ID   string
	Name string
	Desc string
}

func (i tuiSelectItem) Title() string {
	return i.Name
}

func (i tuiSelectItem) Description() string {
	return i.Desc
}

func (i tuiSelectItem) FilterValue() string {
	return i.Name
}

type TUISelectOneParams struct {
	Name        string
	Description string
	Items       []tuiSelectItem
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
	SelectedItem tuiSelectItem
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
