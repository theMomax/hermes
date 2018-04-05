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
    target: hermes
    onSendToQml:
    {
        var data = ""
        if (jsondata != "") {
            data = JSON.parse(jsondata)
        }
        switch(mode) {
        case 0:
            for(var key in data) {
                this.idMap[target][key.toString()] = data[key.toString()]
            }
            break;
        case 1:
            insertElement(target, data, data.template)
            break;
        case 2:
            // open file
            var request = new XMLHttpRequest();
            request.open("GET", data.template, false);
            request.send(null);
            var template = request.responseText;
            insertElement(target, data, template)
            break;
        case 3:
            this.idMap[target].destroy()
            delete this.idMap[target]
            break;
        case 4:
            var request = {}
            for (var i = 0; i < data.properties.length; i++){
                request[data.properties[i]] = this.idMap[target][data.properties[i]]
            }
            hermes.sendToGo(data.eventListener, target, JSON.stringify(request))
            break;
        default:
            // -------------------------------------------
            // Insert your ModeCustom-implementation here.
            // -------------------------------------------
            break;
        }
    }
    function insertElement(target, data, template) {
        // insert variables
        var qmlElement = template.replace(/<\w+>/g, function(match){
            return data.variables[match.replace("<","").replace(">","")]
        })
        // get element id
        var elementId = ""
        qmlElement = qmlElement.replace(/{[^{]*(id\s*:\s*[^\s^;.]+)/i, function(match, p1){
            console.log("id: "+p1)
            elementId = p1.replace(/\s+/g, "").replace("id:","")
            return match
        })
        // create and register element
        this.idMap[elementId] = Qt.createQmlObject(qmlElement, this.idMap[target], data.template)
    }
}