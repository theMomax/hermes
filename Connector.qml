import QtQuick 2.7

Connections
    {
        property var idMap: ({
                                filesColumn:filesColumn,
                                currentCmdTextBinding:currentCmdTextBinding
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
                console.log(target)
                console.log(this.idMap[target])
                for(var key in data) {
                    this.idMap[target].property = key.toString()
                    this.idMap[target].value = data[key.toString()]
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