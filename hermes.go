package hermes

import (
	"bytes"
	"log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/qml"
)

const (
	// ModeSet sets the target's properties according to the provided values.
	// e.g.: target="some_Qml_Text's_id"; jsondata="{"text": "Hello World !",
	// "height": 100}"
	ModeSet int = iota
	// ModeAdd adds a new element specified in the jsondata string to the
	// target element. e.g.: target="some_Qml_Row's_id"; jsondata="
	// {"template": "Text{text: 'Hello World <name> !'}", "variables":
	// {"name": "John Smith"}}". Only the template's root element will be
	// accessible via its id. The id can be hardcoded in the template or
	// dynamically added using the variables section.
	ModeAdd
	// ModeAddFromFile is the same as ModeAdd, but reads the template from a
	// given path. e.g.: target="some_Qml_Row's_id"; jsondata="{"template":
	// "path/to/your/template/from/your/main.qml", "variables":
	// {"name": "John Smith"}}" Only the template's root element will be
	// accessible via its id. The id can be hardcoded in the template or
	// dynamically added using the variables section.
	ModeAddFromFile
	// ModeRemove deletes an element by its qml-id provided in target.
	// The jsondata should be a empty string.
	ModeRemove
)

//go:generate qtmoc

// QmlBridge is the connection between go and qml
type QmlBridge struct {
	core.QObject
	_ func(mode int, target, jsondata string) `signal:"sendToQml"`
	_ func(action, source, jsondata string)   `slot:"sendToGo"`
}

// Controller holds a QmlBridge and the according event-listeners
type Controller struct {
	qmlBridge      *QmlBridge
	eventListeners map[string]func(string, string)
}

// NewBridgeController creates a new controller. This should be done
// before creating a qml window.
func NewBridgeController(engine *qml.QQmlApplicationEngine) *Controller {
	var c = Controller{
		qmlBridge:      NewQmlBridge(nil),
		eventListeners: make(map[string]func(string, string)),
	}
	engine.RootContext().SetContextProperty("hermes", c.qmlBridge)
	c.qmlBridge.ConnectSendToGo(c.interpretQmlCommand)
	return &c
}

// AddEventListener registers a function to be called, when the qml-code
// sends a message with the given event string
func (c *Controller) AddEventListener(event string, action func(string, string)) {
	c.eventListeners[event] = action
}

// RemoveEventListener removes an action, that was previously added
// using AddEventListener()
func (c *Controller) RemoveEventListener(event string) {
	delete(c.eventListeners, event)
}

func (c *Controller) interpretQmlCommand(action, source, jsondata string) {
	log.Println("qml to go: " + string(action) + " | " + source + " | " + jsondata)
	if c.eventListeners[action] != nil {
		c.eventListeners[action](source, jsondata)
	}
}

// SendToQml sends a message to the qml-code. Read the mode-constants'
// comments for further information on how target and jsondata have
// to look like.
func (c *Controller) SendToQml(mode int, target, jsondata string) {
	log.Println("go to qml: " + string(mode) + " | " + target + " | " + jsondata)
	c.qmlBridge.SendToQml(mode, target, jsondata)
}

// SetInQml is shorthand for SendToQml(ModeSet, ...)
func (c *Controller) SetInQml(target, jsondata string) {
	c.SendToQml(ModeSet, target, jsondata)
}

// AddToQml is shorthand for SendToQml(ModeAdd, ...)
func (c *Controller) AddToQml(target, jsondata string) {
	c.SendToQml(ModeAdd, target, jsondata)
}

// AddToQmlFromFilepath is shorthand for SendToQml(ModeAddFromFilepath, ...)
func (c *Controller) AddToQmlFromFilepath(target, jsondata string) {
	c.SendToQml(ModeAddFromFile, target, jsondata)
}

// RemoveFromQml is shorthand for SendToQml(ModeRemove, ...)
func (c *Controller) RemoveFromQml(target string) {
	c.SendToQml(ModeRemove, target, "")
}

// BuildSetModeJSON helps building trivial JSON strings.
// Every odd argument is a property-name, every even one
// the previous property's value
func BuildSetModeJSON(data ...string) string {
	buff := bytes.NewBuffer([]byte{})
	if len(data)%2 != 0 {
		buff.WriteString("{}")
	} else {
		buff.WriteString("{")
		for i, d := range data {
			buff.WriteString(`"`)
			buff.WriteString(d)
			if i%2 == 0 {
				buff.WriteString(`":`)
			} else if i+1 == len(data) {
				buff.WriteString(`"}`)
			} else {
				buff.WriteString(`",`)
			}
		}
	}
	return buff.String()
}

// BuildAddModeJSON helps building trivial JSON strings. The template
// is the template string or filepath. Every odd data argument is a
// variable-name, every even one the previous variable's value
func BuildAddModeJSON(template string, data ...string) string {
	buff := bytes.NewBuffer([]byte{})
	buff.WriteString(`{"template":"`)
	buff.WriteString(template)
	buff.WriteString(`"`)

	if len(data)%2 == 0 {
		buff.WriteString(`, "variables": {`)
		for i, d := range data {
			buff.WriteString(`"`)
			buff.WriteString(d)
			if i%2 == 0 {
				buff.WriteString(`":`)
			} else if i+1 == len(data) {
				buff.WriteString(`"}`)
			} else {
				buff.WriteString(`",`)
			}
		}

	}
	buff.WriteString("}")
	return buff.String()
}
