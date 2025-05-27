package replyme

import (
	"bytes"
	context2 "context"
	"github.com/go-faker/faker/v4"
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestContextToInterface(t *testing.T) {
	ctx := &Context{}

	f := func(c CtxInterface) {}

	f(ctx)
}

func create() *Context {
	command := &Command{
		Name:  "test",
		Usage: faker.Sentence(),
		Subcommands: Commands{
			{
				Name:  "a",
				Usage: faker.Sentence(),
			},
			{
				Name:  "b",
				Usage: faker.Sentence(),
			},
		},
		Flags: Flags{
			&FlagValue[string]{
				Name:           "t1",
				value:          "ofj2o3fj2",
				hasValue:       true,
				preParsedValue: "ofj2o3fj2",
			},
			&FlagValue[bool]{
				Name:           "t2",
				value:          true,
				hasValue:       true,
				preParsedValue: "true",
			},
			&FlagValue[int]{
				Name:           "t3",
				value:          10,
				hasValue:       true,
				preParsedValue: "10",
			},
			&FlagValue[[]string]{
				Name:           "t4",
				value:          []string{"foifof", "foiejrfioejr", "oirfjoer", "eroifj"},
				hasValue:       true,
				preParsedValue: "foifof,foiejrfioejr,oirfjoer,eroifj",
			},
			&FlagValue[[]int]{
				Name:           "t5",
				value:          []int{10, 20, 30, 40},
				hasValue:       true,
				preParsedValue: "10,20,30,40",
			},
		},
		Arguments: []*Argument{
			{
				Name:  "tt2",
				value: "F293fj892u3n2",
			},
		},
	}
	ast := &ASTNode{
		Command:     "test",
		FullCommand: "test --t1=ofj2o3fj2 --t2 --t3=10 F293fj892u3n2",
		CommandTree: []string{
			"test",
		},
		Subcommands: []string{},
		Arguments: []ASTArgument{
			{
				Name:  "tt2",
				Value: "F293fj892u3n2",
			},
		},
		Flags: map[string]map[string][]ASTFlag{
			"test": {
				"t1": {
					{
						Type:  FlagTypeString,
						Value: "ofj2o3fj2",
					},
				},
				"t2": {
					{
						Type:  FlagTypeBool,
						Value: "true",
					},
				},
				"t3": {
					{
						Type:  FlagTypeInt,
						Value: "10",
					},
				},
				"t4": {
					{
						Type:  FlagTypeIntArray,
						Value: "10,20,30,40",
					},
				},
				"t5": {
					{
						Type:  FlagTypeStringArray,
						Value: "foifof,foiejrfioejr,oirfjoer,eroifj",
					},
				},
			},
		},
		Args: []string{
			"F293fj892u3n2",
		},
	}
	ctx, cancel := context2.WithCancel(context2.Background())
	context := &Context{
		ctx:       ctx,
		cancel:    cancel,
		command:   command,
		ast:       ast,
		memory:    new(map[string]interface{}),
		emitLog:   func(msg LogMsg) {},
		stdout:    bytes.NewBuffer(nil),
		stderr:    bytes.NewBuffer(nil),
		startTime: time.Now(),
	}
	context.memory = &map[string]interface{}{}

	return context
}

func TestContext_GetName(t *testing.T) {
	context := create()
	name := context.GetName()

	if name != "test" {
		t.Fatalf("GetName() returns %s, want %s", name, "test")
	}
}

func TestContent_GetCommandNameTree(t *testing.T) {
	context := create()
	tree := context.GetCommandNameTree()

	if !slices.Equal(tree, []string{"test"}) {
		t.Fatalf("GetCommandNameTree() returns %s, want %s", tree, "[]{test}")
	}
}

func TestContent_GetFlagInt(t *testing.T) {
	context := create()
	value := context.GetFlagInt("t3", 0)
	if value != 10 {
		t.Fatalf("GetFlagInt() returns %d, want 10", value)
	}
}

func TestContent_GetFlagString(t *testing.T) {
	context := create()
	value := context.GetFlagString("t1", "")
	if value != "ofj2o3fj2" {
		t.Fatalf("GetFlagString() returns %s, want %s", value, "ofj2o3fj2")
	}
}

func TestContent_GetFlagStringArray(t *testing.T) {
	context := create()
	value := context.GetFlagStringArray("t4")
	if !slices.Equal(value, []string{"foifof", "foiejrfioejr", "oirfjoer", "eroifj"}) {
		t.Fatalf("GetFlagStringArray() returns %s, want %s", value, []string{"foifof", "foiejrfioejr", "oirfjoer", "eroifj"})
	}
}

func TestContent_GetFlagIntArray(t *testing.T) {
	context := create()
	value := context.GetFlagIntArray("t5")
	if !slices.Equal(value, []int{10, 20, 30, 40}) {
		t.Fatalf("GetFlagIntArray() returns %d, want %d", value, []int{10, 20, 30, 40})
	}
}

func TestContent_Print(t *testing.T) {
	ok := false
	context := create()
	context.emitLog = func(msg LogMsg) {
		if msg.Status == LogMsgStatus_Print &&
			msg.Content == "test" {
			ok = true
		}
	}
	context.Print("test")
	if !ok {
		t.Fatalf("Print() returns false, want true")
	}
}

func TestContent_Printf(t *testing.T) {
	ok := false
	context := create()
	context.emitLog = func(msg LogMsg) {
		if msg.Status == LogMsgStatus_Printf &&
			msg.Content == "%s" && msg.Data[0] == "test" {
			ok = true
		}
	}
	context.Printf("%s", "test")
	if !ok {
		t.Fatalf("Printf() returns false, want true")
	}
}

func TestContent_PrintMarkdown(t *testing.T) {
	ok := false
	context := create()
	context.emitLog = func(msg LogMsg) {
		if msg.Status == LogMsgStatus_PrintMarkdown &&
			msg.Content == "t" && msg.Data[0] == "test" {
			ok = true
		}
	}
	context.PrintMarkdown("t", "test")
	if !ok {
		t.Fatalf("PrintMarkdown() returns false, want true")
	}
}

func TestContent_Warn(t *testing.T) {
	ok := false
	context := create()
	context.emitLog = func(msg LogMsg) {
		if msg.Status == LogMsgStatus_Warn &&
			msg.Content == "t" {
			ok = true
		}
	}
	context.Warn("t")
	if !ok {
		t.Fatalf("Warn() returns false, want true")
	}
}

func TestContent_Warnf(t *testing.T) {
	ok := false
	context := create()
	context.emitLog = func(msg LogMsg) {
		if msg.Status == LogMsgStatus_Warnf &&
			msg.Content == "t" && msg.Data[0] == "test" {
			ok = true
		}
	}
	context.Warnf("t", "test")
	if !ok {
		t.Fatalf("Warnf() returns false, want true")
	}
}

func TestContent_Error(t *testing.T) {
	ok := false
	context := create()
	context.emitLog = func(msg LogMsg) {
		if msg.Status == LogMsgStatus_Error &&
			msg.Content == "t" {
			ok = true
		}
	}
	context.Error("t")
	if !ok {
		t.Fatalf("Error() returns false, want true")
	}
}

func TestContent_Errorf(t *testing.T) {
	ok := false
	context := create()
	context.emitLog = func(msg LogMsg) {
		if msg.Status == LogMsgStatus_Errorf &&
			msg.Content == "t" && msg.Data[0] == "test" {
			ok = true
		}
	}
	context.Errorf("t", "test")
	if !ok {
		t.Fatalf("Errorf() returns false, want true")
	}
}

func TestContext_IsCanceled(t *testing.T) {
	context := create()
	if context.IsCancelled() {
		t.Fatalf("IsCancelled() returns true, want false")
	}
	context.cancel()
	if !context.IsCancelled() {
		t.Fatalf("IsCancelled() returns false, want true")
	}
}

func TestContent_SetGet(t *testing.T) {
	context := create()
	context.Set("test1", "test")
	context.Set("test2", 10)
	context.Set("test3", map[string]interface{}{
		"foo": 9383,
		"bar": []string{"foo", "bar"},
	})

	switch {
	case context.Get("test1") != "test":
		t.Fatalf("SetGet() returns %s, want %s", context.Get("test1"), "test")
	case context.Get("test2") != 10:
		t.Fatalf("SetGet() returns %s, want %d", context.Get("test2"), 10)
	case !reflect.DeepEqual(context.Get("test3").(map[string]interface{}), map[string]interface{}{
		"foo": 9383,
		"bar": []string{"foo", "bar"},
	}):
		t.Fatalf("SetGet() returns %v, want %v", context.Get("test3"), map[string]interface{}{
			"foo": 9383,
			"bar": []string{"foo", "bar"},
		})
	}
}

func TestContext_Delete(t *testing.T) {
	context := create()
	context.Set("test1", "test")
	context.Set("test2", 10)
	context.Set("test3", map[string]interface{}{
		"foo": 9383,
		"bar": []string{"foo", "bar"},
	})

	context.Delete("test1")
	context.Delete("test2")
	context.Delete("test3")

	for k, v := range *context.memory {
		t.Fatalf("Memory has data after delete: %v (%v)", k, v)
	}
}

func TestContext_MustGetString(t *testing.T) {
	context := create()
	context.Set("test1", "test")

	if context.MustGetString("test1") != "test" {
		t.Fatalf("MustGetString() returns %s, want %s", context.MustGetString("test1"), "test")
	}
}

func TestContext_MustGetInt(t *testing.T) {
	context := create()
	context.Set("test1", 10)

	if context.MustGetInt("test1") != 10 {
		t.Fatalf("MustGetInt() returns %d, want %d", context.MustGetInt("test1"), 10)
	}
}

func TestContext_Exec(t *testing.T) {
	context := create()
	exec, e, err := context.Exec("echo", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(exec) != "hello" || e != "" {
		t.Fatalf("Exec() returns %s, %s, want %s, %s", exec, e, "hello", "")
	}
}

func TestContent_ExecLive(t *testing.T) {
	logs := make([]LogMsg, 0)
	context := create()
	context.emitLog = func(msg LogMsg) {
		logs = append(logs, msg)
	}
	err := context.ExecLive("echo", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if !(len(logs) == 1 && strings.TrimSpace(logs[0].Content) == "hello") {
		t.Fatalf("ExecLive() returns %v, want %v", logs, "hello")
	}
}

func TestContent_ExecSilient(t *testing.T) {
	context := create()
	err := context.ExecSilent("echo", "hello")
	if err != nil {
		t.Fatal(err)
	}
}
