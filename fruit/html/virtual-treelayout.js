let HORIZONTAL = true;
function init() {
    let $ = go.GraphObject.make;
    myDiagram =
        $(go.Diagram, "myDiagramDiv",
            {
                contentAlignment: go.Spot.Center,
                initialDocumentSpot: go.Spot.Center,
                initialViewportSpot: go.Spot.Center,
                layout:
                    $(VirtualizedTreeLayout,
                        {angle: (HORIZONTAL ? 0 : 90), nodeSpacing: 10}),
                nodeTemplate:
                    $(go.Node, "Horizontal",
                        {
                            isLayoutPositioned: false,
                            toolTip:
                                $("ToolTip",
                                    $(go.TextBlock, {margin: 1},
                                        new go.Binding("text", "",
                                            function (d) {
                                                return "key: " + d.key + "\nbounds: " + d.bounds.toString();
                                            }))
                                )
                        },
                        new go.Binding("position", "bounds", function (b) {
                            console.log(b);
                            return b.position;
                        }).makeTwoWay(function (p, d) {
                            return new go.Rect(p.x, p.y, d.bounds.width, d.bounds.height);
                        }),
                        $(go.Panel, "Auto",
                            $(go.Shape, "Rectangle",
                                new go.Binding("fill", "color")),
                            $(go.TextBlock,
                                {margin: 1, editable: true, textAlign: "left"},
                                new go.Binding("text", "text")),
                            $(go.Panel, "Auto",
                                {
                                    alignment: go.Spot.Right,
                                    portId: "from",
                                    fromLinkable: true,
                                    cursor: "pointer",
                                    click: addNodeAndLink
                                }
                            )
                        ),
                        $("TreeExpanderButton")
                    ),
                // Define the template for Links
                linkTemplate:
                    $(go.Link,
                        {
                            fromSpot: (HORIZONTAL ? go.Spot.Right : go.Spot.Bottom),
                            toSpot: (HORIZONTAL ? go.Spot.Left : go.Spot.Top)
                        },
                        $(go.Shape),
                        //$(go.Shape, { toArrow: "OpenTriangle", stroke: "black"})
                    ),
                "SelectionMoved":
                    function (e) {
                        e.subject.each(function (n) {
                            if (n instanceof go.Node) n.data.points = undefined;
                        })
                    },
                "animationManager.isEnabled":
                    false
            });


    // This model includes all of the data
    myWholeModel = $(go.TreeModel);

    // The virtualized layout works on the full model, not on the Diagram Nodes and Links
    myDiagram.layout.model = myWholeModel;

    function addNodeAndLink(e, obj) {
        let adornment = obj.part;
        let diagram = e.diagram;
        let model = diagram.model;
        diagram.startTransaction("Add State");

        // get the node data for which the user clicked the button
        let fromNode = adornment.adornedPart;
        let fromData = fromNode.data;
        let id = myWholeModel.nodeDataArray.length + 1;
        // create a new "State" data object, positioned off to the right of the adorned Node
        let p = fromData.bounds.copy();
        p.x += 150;
        let toData = {key: id, text: "NewData", parent: fromData.key, bounds: p, color: "#ffc107"};
        model.addNodeData(toData);
        myWholeModel.nodeDataArray.push(toData);
        // select the new Node
        let newnode = diagram.findNodeForData(toData);
        diagram.select(newnode);
        diagram.commitTransaction("Add State");
    }

    myDiagram.nodeTemplate.selectionAdornmentTemplate =
        $(go.Adornment, "Spot",
            $(go.Panel, "Auto",
                $(go.Shape, "RoundedRectangle",
                    {fill: null}),
                $(go.Placeholder)
            ),
            // the button to create a "next" node, at the top-right corner
            $("Button",
                {
                    alignment: go.Spot.TopRight,
                    click: addNodeAndLink  // this function is defined below
                },
                $(go.Shape, "PlusLine", {width: 6, height: 6})
            ) // end button
        ); // end Adornment

    // Do not set myDiagram.model = myWholeModel -- that would create a zillion Nodes and Links!
    // In the future Diagram may have built-in support for virtualization.
    // For now, we have to implement virtualization ourselves by having the Diagram's model
    // be different than the "real" model.
    myDiagram.model =   // this only holds nodes that should be in the viewport
        $(go.TreeModel);  // must match the model, above

    // for now, we have to implement virtualization ourselves
    myDiagram.isVirtualized = true;
    myDiagram.addDiagramListener("ViewportBoundsChanged", onViewportChanged);

    // This is a status message
    myLoading =
        $(go.Part,  // this has to set the location or position explicitly
            {location: new go.Point(0, 0)},
            $(go.TextBlock, "Wait a few minutes",
                {stroke: "black", font: "20pt sans-serif"}));

    // temporarily add the status indicator
    myDiagram.add(myLoading);

    // Allow the myLoading indicator to be shown now,
    // but allow objects added in load to also be considered part of the initial Diagram.
    // If you are not going to add temporary initial Parts, don't call delayInitialization.
    myDiagram.delayInitialization(load);
}

function getColor(color) {
    switch (color) {
        case "primary":
            return "#52ce60";
        case "running":
            return "#fd7e14";
        case "archive":
            return "#e83e8c";
        default:
            return "#ffffff"
    }
}

function load() {
    // create a lot of data for the myWholeModel
    axios.get('http://localhost:8080/worker/fruit/api/v1', {
        params: {
            type: "tracing",
            time: "2019"
        }
    }).then(function (response) {
        let treedata = [];
        console.log(response.data);
        for (let i = 0; i < response.data.nodeDataArray.length; i++) {
            let t = response.data.nodeDataArray[i];
            let c = getColor(t.color);
            let d = {
                key: t.key,
                color: c,
                parent: t.parent,
                text: t.text
            };
            //!!!???@@@ this needs to be customized to account for your chosen Node template
            d.bounds = new go.Rect(0, 0, 100, 20);
            treedata.push(d);
        }
        myWholeModel.nodeDataArray = treedata;

        // remove the status indicator
        myDiagram.remove(myLoading);

    })
        .catch(function (error) {
            console.log(error);
        })
        .finally(function () {

        });


}


// The following functions implement virtualization of the Diagram
// Assume data.bounds is a Rect of the area occupied by the Node in document coordinates.

// The normal mechanism for determining the size of the document depends on all of the
// Nodes and Links existing, so we need to use a function that depends only on the model data.
function computeDocumentBounds(model) {
    let b = new go.Rect();
    let ndata = model.nodeDataArray;
    for (let i = 0; i < ndata.length; i++) {
        let d = ndata[i];
        if (!d.bounds) continue;
        if (i === 0) {
            b.set(d.bounds);
        } else {
            b.unionRect(d.bounds);
        }
    }
    return b;
}

// As the user scrolls or zooms, make sure the Parts (Nodes and Links) exist in the viewport.
function onViewportChanged(e) {
    let diagram = e.diagram;
    // make sure there are Nodes for each node data that is in the viewport
    // or that is connected to such a Node
    let viewb = diagram.viewportBounds;  // the new viewportBounds
    let model = diagram.model;

    let oldskips = diagram.skipsUndoManager;
    diagram.skipsUndoManager = true;

    let b = new go.Rect();
    let ndata = myWholeModel.nodeDataArray;
    for (let i = 0; i < ndata.length; i++) {
        let n = ndata[i];
        if (!n.bounds) continue;
        if (n.bounds.intersectsRect(viewb)) {
            model.addNodeData(n);
        }
        if (model instanceof go.TreeModel) {
            // make sure links to all parent nodes appear
            let parentkey = myWholeModel.getParentKeyForNodeData(n);
            let parent = myWholeModel.findNodeDataForKey(parentkey);
            if (parent !== null) {
                if (n.bounds.intersectsRect(viewb)) {  // N is inside viewport
                    model.addNodeData(parent);  // so that link to parent appears
                    let child = diagram.findNodeForData(n);
                    if (child !== null) {
                        let link = child.findTreeParentLink();
                        if (link !== null) {
                            // do this now to avoid delayed routing outside of transaction
                            link.fromNode.ensureBounds();
                            link.toNode.ensureBounds();
                            if (child.data.fromSpot) link.fromSpot = child.data.fromSpot;
                            if (child.data.toSpot) link.toSpot = child.data.toSpot;
                            if (child.data.points) {
                                link.points = child.data.points;
                            } else {
                                link.updateRoute();
                            }
                        }
                    }
                } else {  // N is outside of viewport
                    // see if there's a parent that is in the viewport,
                    // or if the link might cross over the viewport
                    b.set(n.bounds);
                    b.unionRect(parent.bounds);
                    if (b.intersectsRect(viewb)) {
                        model.addNodeData(n);  // add N so that link to parent appears
                        let child = diagram.findNodeForData(n);
                        if (child !== null) {
                            let link = child.findTreeParentLink();
                            if (link !== null) {
                                // do this now to avoid delayed routing outside of transaction
                                link.fromNode.ensureBounds();
                                link.toNode.ensureBounds();
                                if (child.data.fromSpot) link.fromSpot = child.data.fromSpot;
                                if (child.data.toSpot) link.toSpot = child.data.toSpot;
                                if (child.data.points) {
                                    link.points = child.data.points;
                                } else {
                                    link.updateRoute();
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    if (model instanceof go.GraphLinksModel) {
        let ldata = myWholeModel.linkDataArray;
        for (let i = 0; i < ldata.length; i++) {
            let l = ldata[i];
            let fromkey = myWholeModel.getFromKeyForLinkData(l);
            if (fromkey === undefined) continue;
            let from = myWholeModel.findNodeDataForKey(fromkey);
            if (from === null || !from.bounds) continue;

            let tokey = myWholeModel.getToKeyForLinkData(l);
            if (tokey === undefined) continue;
            let to = myWholeModel.findNodeDataForKey(tokey);
            if (to === null || !to.bounds) continue;

            b.set(from.bounds);
            b.unionRect(to.bounds);
            if (b.intersectsRect(viewb)) {
                // also make sure both connected nodes are present,
                // so that link routing is authentic
                model.addNodeData(from);
                model.addNodeData(to);
                model.addLinkData(l);
                let link = diagram.findLinkForData(l);
                if (link !== null) {
                    // do this now to avoid delayed routing outside of transaction
                    link.fromNode.ensureBounds();
                    link.toNode.ensureBounds();
                    if (link.data.fromSpot) link.fromSpot = link.data.fromSpot;
                    if (link.data.toSpot) link.toSpot = link.data.toSpot;
                    //if (link.data.points) {
                    //  link.points = link.data.points;
                    //} else {
                    link.updateRoute();
                    //}
                }
            }
        }
    }

    diagram.skipsUndoManager = oldskips;

    if (myRemoveTimer === null) {
        // only remove offscreen nodes after a delay
        myRemoveTimer = setTimeout(function () {
            removeOffscreen(diagram);
        }, 3000);
    }

    updateCounts();  // only for this sample
}

// occasionally remove Parts that are offscreen from the Diagram
let myRemoveTimer = null;

function removeOffscreen(diagram) {
    myRemoveTimer = null;

    let viewb = diagram.viewportBounds;
    let model = diagram.model;
    let remove = [];  // collect for later removal
    let removeLinks = new go.Set();  // links connected to a node data to remove
    let it = diagram.nodes;
    while (it.next()) {
        let n = it.value;
        let d = n.data;
        if (d === null) continue;
        if (!n.actualBounds.intersectsRect(viewb) && !n.isSelected) {
            // even if the node is out of the viewport, keep it if it is selected or
            // if any link connecting with the node is still in the viewport
            if (!n.linksConnected.any(function (l) {
                return l.actualBounds.intersectsRect(viewb);
            })) {
                remove.push(d);
                if (model instanceof go.GraphLinksModel) {
                    removeLinks.addAll(n.linksConnected);
                }
            }
        }
    }

    if (remove.length > 0) {
        let oldskips = diagram.skipsUndoManager;
        diagram.skipsUndoManager = true;
        model.removeNodeDataCollection(remove);
        if (model instanceof go.GraphLinksModel) {
            removeLinks.each(function (l) {
                if (!l.isSelected) model.removeLinkData(l.data);
            });
        }
        diagram.skipsUndoManager = oldskips;
    }

    updateCounts();
}

function VirtualizedTreeLayout() {
    go.TreeLayout.call(this);
    this.isOngoing = false;
    this.model = null;
}

go.Diagram.inherit(VirtualizedTreeLayout, go.TreeLayout);

VirtualizedTreeLayout.prototype.createNetwork = function () {
    return new VirtualizedTreeNetwork(this);  // defined below
};


VirtualizedTreeLayout.prototype.makeNetwork = function (coll) {
    let net = this.createNetwork();
    net.addData(this.model);
    return net;
};

VirtualizedTreeLayout.prototype.commitLayout = function () {
    VirtualizedTreeEdge._dummyLink = this.diagram.linkTemplate.copy();
    VirtualizedTreeEdge._dummyFromNode = this.diagram.nodeTemplate.copy();
    VirtualizedTreeEdge._dummyToNode = this.diagram.nodeTemplate.copy();
    VirtualizedTreeEdge._dummyLink.fromNode = VirtualizedTreeEdge._dummyFromNode;
    VirtualizedTreeEdge._dummyLink.toNode = VirtualizedTreeEdge._dummyToNode;
    this.diagram.add(VirtualizedTreeEdge._dummyFromNode);
    this.diagram.add(VirtualizedTreeEdge._dummyToNode);
    this.diagram.add(VirtualizedTreeEdge._dummyLink);

    go.TreeLayout.prototype.commitLayout.call(this);
    // can't depend on regular bounds computation that depends on all Nodes existing
    this.diagram.fixedBounds = computeDocumentBounds(this.model);
    // update the positions of any existing Nodes
    this.diagram.nodes.each(function (node) {
        node.updateTargetBindings();
    });

    this.diagram.remove(VirtualizedTreeEdge._dummyFromNode);
    this.diagram.remove(VirtualizedTreeEdge._dummyToNode);
    this.diagram.remove(VirtualizedTreeEdge._dummyLink);
};

VirtualizedTreeLayout._cachedLinks = [];

VirtualizedTreeLayout.prototype.setPortSpots = function (v) {
    v.destinationEdges.each(function (e) {
        e.link = VirtualizedTreeLayout._cachedLinks.pop() || new go.Link();
    });
    go.TreeLayout.prototype.setPortSpots.call(this, v);
    v.destinationEdges.each(function (e) {
        if (e.data) {
            e.data.fromSpot = e.link.fromSpot.copy();
            e.data.toSpot = e.link.toSpot.copy();
        }
        VirtualizedTreeLayout._cachedLinks.push(e.link);
        e.link = null;
    });
};


function VirtualizedTreeNetwork(layout) {
    go.TreeNetwork.call(this, layout);
}

go.Diagram.inherit(VirtualizedTreeNetwork, go.TreeNetwork);

VirtualizedTreeNetwork.prototype.createEdge = function () {
    return new VirtualizedTreeEdge(this);
};

VirtualizedTreeNetwork.prototype.addData = function (model) {
    if (model instanceof go.TreeModel) {
        let dataVertexMap = new go.Map();
        let ndata = model.nodeDataArray;
        for (let i = 0; i < ndata.length; i++) {
            let d = ndata[i];
            let v = this.createVertex();
            v.data = d;  // associate this Vertex with data, not a Node
            dataVertexMap.set(d, v);
            this.addVertex(v);
        }

        for (let i = 0; i < ndata.length; i++) {
            let child = ndata[i];
            let parentkey = model.getParentKeyForNodeData(child);
            let parent = model.findNodeDataForKey(parentkey);
            if (parent !== null) {  // if there is a parent, there should be an edge
                // now find corresponding vertexes
                let f = dataVertexMap.get(parent);
                let t = dataVertexMap.get(child);
                if (f === null || t === null) continue;  // skip
                // create and add VirtualizedTreeEdge
                let e = this.createEdge();
                e.data = child;  // associate this Edge with data, not a Link
                e.fromVertex = f;
                e.toVertex = t;
                this.addEdge(e);
            }
        }
    } else {
        throw new Error("can only handle TreeModel data");
    }
};

// end VirtualizedTreeNetwork class

function VirtualizedTreeEdge(network) {
    go.TreeEdge.call(this, network);
}

go.Diagram.inherit(VirtualizedTreeEdge, go.TreeEdge);

VirtualizedTreeEdge._dummyLink = null;
VirtualizedTreeEdge._dummyFromNode = null;
VirtualizedTreeEdge._dummyToNode = null;

VirtualizedTreeEdge.prototype.commit = function () {
    let parentv = this.fromVertex;
    if (!parentv) return;
    let routed = (parentv.alignment === go.TreeLayout.AlignmentStart || parentv.alignment === go.TreeLayout.AlignmentEnd);
    if (this.data && routed) {
        this.link = VirtualizedTreeEdge._dummyLink;
        this.link.fromNode.position = new go.Point(this.fromVertex.x, this.fromVertex.y);
        this.link.toNode.position = new go.Point(this.toVertex.x, this.toVertex.y);
        this.link.fromNode.ensureBounds();
        this.link.toNode.ensureBounds();
        this.link.updateRoute();
    }
    go.TreeEdge.prototype.commit.call(this);
    if (this.data && routed) {
        this.data.points = this.link.points.copy();
        if (this.link.fromNode.actualBounds.x < this.link.toNode.actualBounds.x - 150) {
            console.log(this.link.points)
        }
        this.link = null;
    }
};

// end of VirtualizedTree[Layout/Network] classes

// This function is only used in this sample to demonstrate the effects of the virtualization.
// In a real application you would delete this function and all calls to it.
function updateCounts() {
    document.getElementById("nodes").textContent = myWholeModel.nodeDataArray.length.toString();
    document.getElementById("node-shown").textContent = myDiagram.nodes.count;
    document.getElementById("link-shown").textContent = myDiagram.links.count;
}

function save() {
    document.getElementById("mySavedModel").value = myDiagram.model.toJson();
}
