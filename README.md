# hermes

A high-level communication interface for github.com/therecipe/qt. This package allows you to dynamically add Qml-Elements (from Qml-code as a string or file), set and read their properties and delete them.

## Get Started

Go get this repository:
````
go get github.com/hoffx/hermes
````
Integrate `hermes` into your Qml application on Go side by importing it and creating a new `hermes.Controller`:
````Go
// main.go
// Create application
app := gui.NewQGuiApplication(len(os.Args), os.Args)

// Create a QML application engine
engine := qml.NewQQmlApplicationEngine(nil)

hController = hermes.NewBridgeController(engine)
hController.AddEventListener("someEvent", someFunctionToCallAtEvent)
hController.AddEventListener("someEvent2", someFunctionToCallAtEvent2)

// Load the main qml file
window := qml.NewQQmlComponent5(engine, core.NewQUrl3("qml/main.qml", 0), nil)
root = window.Create(engine.RootContext())

// Execute app
gui.QGuiApplication_Exec()
````
Integrate `hermes` into your Qml application on Qml side by adding the code-snippet from `hermes.qml` to the ApplicationWindow element of your `main.qml` file:
````Qml
import QtQuick 2.7
ApplicationWindow {
    id: window
    visible: true
    ...

    Connections
    {
        property var idMap: ({
            // -------------------------------------------
            // List all the ids here, you have to access
            // from Go. This is needed for
            // string to object mapping.
            // -------------------------------------------
            //window:window,
            //yourelement1:yourelement1
        })
        ...
    }
````

## How to use it

On Go side it's easy. Just use your `hController`'s functions. They aren't documented that bad... But make sure every qml element you want to access has been added to the idMap - manually by you (for the static part) or by the `AddToQml` function. Just read the documentation...


In Qml, use the `hermes.sendToGo()` command:
````Qml
Button{
    id: yourButton
    onClicked: hermes.sendToGo("someEvent", "yourButton", '{ "extra_information_property": "extra_information_value" }')
}
````
