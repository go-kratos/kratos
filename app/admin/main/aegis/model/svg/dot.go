package svg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/log"
)

var tpl = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>流程网描述图</title>

        <script src="https://cdn.bootcss.com/viz.js/2.1.1/viz.js"></script>
        <script src="https://cdn.bootcss.com/viz.js/2.1.1/full.render.js"></script>

        <script>
            var json={{.ActionDesc}};
            var actions = {{.Actions}};
            window.onload = function () {
                    for( let i =0; i<actions.length; i++){
                        document.getElementById(actions[i]).onclick = function(e){
                          var x = e.clientX;
                          var y = e.clientY;

                          var content = document.getElementById("content");
                          
                          var table = '<table bgcolor="LightSteelBlue" border="1">'
                          var data = JSON.parse(json[actions[i]]);
                          for (var key in data){
                            table +='<tr>'
                            table+='<td>'+key+'</td>'
                            table+='<td>'+data[key]+'</td>'
                            table +='</tr>'
                          }
                          table += '</table>'

                          content.innerHTML = table;
                          content.style.left = x+"px";
                          content.style.top = y+"px";
                          content.style.position="absolute";
                          content.style.display = "block";
                          content.disabled = true
                          console.log(content);
                        }
                    };
            };
        </script>
    </head>
    <body>
        <label id="content" background-color="red"></label>
        <script>
            var viz = new Viz();
            
            viz.renderSVGElement('DotCode')
            .then(function(element) {
              document.body.appendChild(element);
            })
            .catch(error => {
              // Create a new Viz instance (@see Caveats page for more info)
              viz = new Viz();
          
              // Possibly display the error
              console.error(error);
            });
          </script>
    </body>
</html>
`

var (
	_labelFlow       = `%s [label="%s",shape=circle]`
	_labelTransition = `%s [label="%s",shape=box,color=lightgrey,style=filled;]`
	_labelToken      = `%d [label="%s",shape=Mdiamond]`
	_labelDirection  = `%s->%s`
)

//NetView .
type NetView struct {
	*template.Template
	Dot  *Dot
	Data struct {
		NetDesc    string
		DotCode    string
		Actions    []string
		ActionDesc map[string]string
	}
}

//NewNetView .
func NewNetView() *NetView {
	return &NetView{}
}

//SetDot .
func (nv *NetView) SetDot(dot *Dot) {
	nv.Dot = dot
	nv.Data.DotCode = dot.String()
	nv.Data.Actions = append(dot.nodes[:], dot.edges...)
	nv.Data.ActionDesc = dot.mapactions

	ntpl := strings.Replace(tpl, "DotCode", nv.Data.DotCode, 1)
	nv.Template = template.New("流程网概览")
	nv.Template, _ = nv.Template.Parse(ntpl)
}

/*
节点用record表示，节点绑定的令牌用在同一个record里面
变迁用diamond表示，若变迁绑定了令牌，则用图包含变迁和令牌
*/
//Dot .
type Dot struct {
	*bytes.Buffer
	nodes      []string
	edges      []string
	mapactions map[string]string
	flows      map[int64]*net.Flow
	trans      map[int64]*net.Transition
	dirs       map[int64]*net.Direction
	tokens     map[string][]*net.TokenBindDetail
}

//NewDot .
func NewDot() *Dot {
	d := &Dot{
		mapactions: make(map[string]string),
		flows:      make(map[int64]*net.Flow),
		trans:      make(map[int64]*net.Transition),
		dirs:       make(map[int64]*net.Direction),
		tokens:     make(map[string][]*net.TokenBindDetail),
	}
	return d
}

//StartDot .
func (d *Dot) StartDot() *Dot {
	d.Buffer = bytes.NewBufferString(`digraph net{rankdir=LR;`)
	return d
}

//End .
func (d *Dot) End() string {
	d.WriteString("}")
	return d.String()
}

//AddTokenBinds .
func (d *Dot) AddTokenBinds(tbs ...*net.TokenBindDetail) *Dot {
	for _, tb := range tbs {
		key := fmt.Sprintf("%d_%d", tb.Type, tb.ElementID)
		if tks, ok := d.tokens[key]; ok {
			d.tokens[key] = append(tks, tb)
		} else {
			d.tokens[key] = []*net.TokenBindDetail{tb}
		}
	}
	return d
}

//AddFlow .
func (d *Dot) AddFlow(flows ...*net.Flow) *Dot {
	for _, flow := range flows {
		d.flows[flow.ID] = flow
		node := fmt.Sprintf(_labelFlow, flow.Name, flow.ChName)
		//
		nodeid := fmt.Sprintf("node%d", len(d.nodes)+1)
		d.nodes = append(d.nodes, nodeid)
		bs, _ := json.Marshal(flow)
		d.mapactions[nodeid] = string(bs)

		//便利token，查找绑定
		if tks, ok := d.tokens[fmt.Sprintf("1_%d", flow.ID)]; ok {
			d.WriteString("subgraph cluster_" + flow.Name + " {")
			d.WriteString(node + ";")
			for _, tk := range tks {
				node := fmt.Sprintf(_labelToken, tk.ID, tk.ChName)
				d.WriteString(node + ";")
				nodeid := fmt.Sprintf("node%d", len(d.nodes)+1)
				d.nodes = append(d.nodes, nodeid)
				bs, _ := json.Marshal(tk)
				d.mapactions[nodeid] = string(bs)
			}
			d.WriteString("}")
		} else {
			d.WriteString(node + ";")
		}
	}
	return d
}

//AddTransitions .
func (d *Dot) AddTransitions(trans ...*net.Transition) *Dot {
	for _, tran := range trans {
		d.trans[tran.ID] = tran
		node := fmt.Sprintf(_labelTransition, tran.Name, tran.ChName)
		nodeid := fmt.Sprintf("node%d", len(d.nodes)+1)
		d.nodes = append(d.nodes, nodeid)
		bs, _ := json.Marshal(tran)
		d.mapactions[nodeid] = string(bs)

		//便利token，查找绑定
		if tks, ok := d.tokens[fmt.Sprintf("2_%d", tran.ID)]; ok {
			d.WriteString("subgraph cluster_" + tran.Name + " {")
			d.WriteString(node + ";")
			for _, tk := range tks {
				node := fmt.Sprintf(_labelToken, tk.ID, tk.ChName)
				d.WriteString(node + ";")
				nodeid := fmt.Sprintf("node%d", len(d.nodes)+1)
				d.nodes = append(d.nodes, nodeid)
				bs, _ := json.Marshal(tk)
				d.mapactions[nodeid] = string(bs)
			}
			d.WriteString("}")
		} else {
			d.WriteString(node + ";")
		}
	}
	return d
}

//AddDirections .
func (d *Dot) AddDirections(dirs ...*net.Direction) *Dot {
	for _, dir := range dirs {
		var (
			start, end string
			flow       *net.Flow
			trans      *net.Transition
		)

		flow = d.flows[dir.FlowID]
		trans = d.trans[dir.TransitionID]

		if flow == nil || trans == nil {
			log.Error("invalid direction(%+v)", dir)
			continue
		}

		if dir.Direction == 1 {
			start, end = flow.Name, trans.Name
		}
		if dir.Direction == 2 {
			start, end = trans.Name, flow.Name
		}

		edge := fmt.Sprintf(_labelDirection, start, end)
		d.WriteString(edge + ";")
		edgeid := fmt.Sprintf("edge%d", len(d.edges)+1)
		bs, _ := json.Marshal(dir)
		d.mapactions[edgeid] = string(bs)
		d.edges = append(d.edges, edgeid)
	}
	return d
}

var (
	flow1 = &net.Flow{ID: 1, ChName: "节点1", Name: "flow1"}
	flow2 = &net.Flow{ID: 2, ChName: "节点2", Name: "flow2"}
	flow3 = &net.Flow{ID: 3, ChName: "节点3", Name: "flow3"}
	flow4 = &net.Flow{ID: 4, ChName: "节点4", Name: "flow4"}
	flow5 = &net.Flow{ID: 5, ChName: "节点5", Name: "flow5"}
	flow6 = &net.Flow{ID: 6, ChName: "节点6", Name: "flow6"}
	flow7 = &net.Flow{ID: 7, ChName: "节点7", Name: "flow7"}
	tran1 = &net.Transition{ID: 1, ChName: "变迁1", Name: "tran1"}
	tran2 = &net.Transition{ID: 2, ChName: "变迁2", Name: "tran2"}
	tran3 = &net.Transition{ID: 3, ChName: "变迁3", Name: "tran3"}
	dir1  = &net.Direction{ID: 1, FlowID: 1, TransitionID: 1, Direction: 1}
	dir2  = &net.Direction{ID: 2, FlowID: 2, TransitionID: 1, Direction: 2}
	dir3  = &net.Direction{ID: 3, FlowID: 2, TransitionID: 2, Direction: 1}
	dir4  = &net.Direction{ID: 4, FlowID: 2, TransitionID: 3, Direction: 1}
	dir5  = &net.Direction{ID: 5, FlowID: 3, TransitionID: 2, Direction: 2}
	dir6  = &net.Direction{ID: 6, FlowID: 4, TransitionID: 2, Direction: 2}
	dir7  = &net.Direction{ID: 7, FlowID: 5, TransitionID: 3, Direction: 2}
	dir8  = &net.Direction{ID: 8, FlowID: 6, TransitionID: 2, Direction: 2}
	dir9  = &net.Direction{ID: 9, FlowID: 7, TransitionID: 3, Direction: 2}

	tk1 = &net.TokenBindDetail{
		ID:        1,
		Type:      1,
		ElementID: 1,
		ChName:    "待审核",
	}
	tk2 = &net.TokenBindDetail{
		ID:        2,
		Type:      2,
		ElementID: 1,
		ChName:    "通过",
	}
)

func DebugSVG() (nv *NetView) {
	dot := NewDot()
	dot.StartDot().AddTokenBinds(tk1, tk2).
		AddFlow(flow1, flow2, flow3, flow4, flow5, flow6, flow7).
		AddTransitions(tran1, tran2, tran3).
		AddDirections(dir1, dir2, dir3, dir4, dir5, dir6, dir7, dir8, dir9).
		End()
	nv = NewNetView()
	nv.SetDot(dot)
	return nv
}
