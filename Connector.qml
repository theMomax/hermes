import QtQuick 2.7

Connections
{
    property var idMap: ({
                            filesColumn:filesColumn,
                            currentCmdTextBinding:currentCmdTextBinding
                        })
    target: qmlBridge
    onSendToQml:
    {
        var data = ""
        if (jsondata != "") {
            data = JSON.parse(jsondata)
        }
        switch(mode) {
        case 0:
            for(var key in data) {
                this.idMap[target].property = key.toString()
                this.idMap[target].value = data[key.toString()]
            }
            break;
        case 1:
            var qmlElement = data.template.replace(/<\w+>/i, function(match){
                return data.variables[match.replace("<","").replace(">","")]
            })
            Qt.createQmlObject(qmlElement, this.idMap[target], target + "_ChildTemplate")
            break;
        case 2:
            // open file
            var request = new XMLHttpRequest();
            request.open("GET", data.template, false);
            request.send(null);
            var template = request.responseText;
            var qmlElement = template.replace(/<\w+>/i, function(match){
                return data.variables[match.replace("<","").replace(">","")]
            })
            Qt.createQmlObject(qmlElement, this.idMap[target], data.template)
            break;
        case 3:
            this.idMap[target].destroy()
            break;
        }
    }
}