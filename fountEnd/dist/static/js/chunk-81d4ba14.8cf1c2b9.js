(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-81d4ba14"],{"22ce":function(e,t,r){"use strict";r.d(t,"f",(function(){return i})),r.d(t,"d",(function(){return s})),r.d(t,"e",(function(){return l})),r.d(t,"b",(function(){return o})),r.d(t,"i",(function(){return n})),r.d(t,"c",(function(){return c})),r.d(t,"j",(function(){return m})),r.d(t,"a",(function(){return u})),r.d(t,"h",(function(){return d})),r.d(t,"g",(function(){return p}));var a=r("b775");function i(e){return Object(a["a"])({url:"/service/serviceList",method:"get",params:e})}function s(e){return Object(a["a"])({url:"/service/serviceDelete",method:"get",params:e})}function l(e){return Object(a["a"])({url:"/service/serviceDetail",method:"get",params:e})}function o(e){return Object(a["a"])({url:"/service/serviceAddHttp",method:"post",params:e})}function n(e){return Object(a["a"])({url:"/service/serviceUpdateHttp",method:"post",params:e})}function c(e){return Object(a["a"])({url:"/service/serviceAddTcp",method:"post",params:e})}function m(e){return Object(a["a"])({url:"/service/serviceUpdateTcp",method:"post",params:e})}function u(e){return Object(a["a"])({url:"/service/serviceAddGrpc",method:"post",params:e})}function d(e){return Object(a["a"])({url:"/service/serviceUpdateGrpc",method:"post",params:e})}function p(e){return Object(a["a"])({url:"/service/serviceStat",method:"get",params:e})}},7207:function(e,t,r){"use strict";r.r(t);var a=function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("div",{staticClass:"mixin-components-container"},[r("el-row",[r("el-card",{staticClass:"box-card"},[r("div",{staticClass:"clearfix",attrs:{slot:"header"},slot:"header"},[!1===e.isEdit?r("span",[e._v("创建HTTP服务")]):e._e(),!0===e.isEdit?r("span",[e._v("修改HTTP服务")]):e._e()]),r("div",{staticStyle:{"margin-bottom":"50px"}},[r("el-form",{ref:"form",attrs:{model:e.form,"label-width":"140px"}},[r("el-form-item",{staticClass:"is-required",attrs:{label:"服务名称"}},[r("el-input",{attrs:{placeholder:"6-128位字母数字下划线",disabled:!0===e.isEdit},model:{value:e.form.serviceName,callback:function(t){e.$set(e.form,"serviceName",t)},expression:"form.serviceName"}})],1),r("el-form-item",{staticClass:"is-required",attrs:{label:"服务描述"}},[r("el-input",{attrs:{placeholder:"最多255字符，必填"},model:{value:e.form.serviceDesc,callback:function(t){e.$set(e.form,"serviceDesc",t)},expression:"form.serviceDesc"}})],1),r("el-form-item",{staticClass:"is-required",attrs:{label:"接入类型"}},[r("el-input",{staticClass:"input-with-select",attrs:{disabled:!0===e.isEdit,placeholder:"路径格式：/user/,域名格式：www.test.com"},model:{value:e.form.rule,callback:function(t){e.$set(e.form,"rule",t)},expression:"form.rule"}},[r("el-select",{staticStyle:{width:"80px"},attrs:{slot:"prepend",placeholder:"请选择",disabled:!0===e.isEdit},slot:"prepend",model:{value:e.form.ruleType,callback:function(t){e.$set(e.form,"ruleType",t)},expression:"form.ruleType"}},[r("el-option",{attrs:{label:"路径",value:0}}),r("el-option",{attrs:{label:"域名",value:1}})],1)],1)],1),r("el-form-item",{attrs:{label:"支持https"}},[r("el-switch",{attrs:{"active-value":1,"inactive-value":0},model:{value:e.form.needHttps,callback:function(t){e.$set(e.form,"needHttps",t)},expression:"form.needHttps"}}),r("span",{staticStyle:{color:"#606266","font-weight":"700","margin-left":"20px"}},[e._v("支持strip_uri")]),r("el-switch",{staticStyle:{"margin-left":"10px"},attrs:{"active-value":1,"inactive-value":0},model:{value:e.form.needStripUrl,callback:function(t){e.$set(e.form,"needStripUrl",t)},expression:"form.needStripUrl"}}),r("span",{staticStyle:{color:"#606266","font-weight":"700","margin-left":"20px"}},[e._v("支持websocket")]),r("el-switch",{staticStyle:{"margin-left":"10px"},attrs:{"active-value":1,"inactive-value":0},model:{value:e.form.needWebsocket,callback:function(t){e.$set(e.form,"needWebsocket",t)},expression:"form.needWebsocket"}})],1),r("el-form-item",{attrs:{label:"URL重写"}},[r("el-input",{attrs:{type:"textarea",autosize:"",placeholder:"格式：^/gateway/test_service(.*) $1 多条换行"},model:{value:e.form.urlRewrite,callback:function(t){e.$set(e.form,"urlRewrite",t)},expression:"form.urlRewrite"}})],1),r("el-form-item",{attrs:{label:"Header转换"}},[r("el-input",{attrs:{type:"textarea",autosize:"",placeholder:"header转换支持 add(增加)/del(删除)/edit(修改) 格式：add headerName headValue"},model:{value:e.form.headerTransfor,callback:function(t){e.$set(e.form,"headerTransfor",t)},expression:"form.headerTransfor"}})],1),r("el-form-item",{attrs:{label:"开启验证"}},[r("el-switch",{attrs:{"active-value":1,"inactive-value":0},model:{value:e.form.openAuth,callback:function(t){e.$set(e.form,"openAuth",t)},expression:"form.openAuth"}})],1),r("el-form-item",{attrs:{label:"IP白名单"}},[r("el-input",{attrs:{type:"textarea",autosize:"",placeholder:"格式：127.0.0.1 多条换行，白名单优先于黑名单"},model:{value:e.form.whiteList,callback:function(t){e.$set(e.form,"whiteList",t)},expression:"form.whiteList"}})],1),r("el-form-item",{attrs:{label:"IP黑名单"}},[r("el-input",{attrs:{type:"textarea",autosize:"",placeholder:"格式：127.0.0.1 多条换行，白名单优先于黑名单"},model:{value:e.form.blackList,callback:function(t){e.$set(e.form,"blackList",t)},expression:"form.blackList"}})],1),r("el-form-item",{attrs:{label:"客户端限流"}},[r("el-input",{attrs:{placeholder:"0表示无限制"},model:{value:e.form.clientIpFlowLimit,callback:function(t){e.$set(e.form,"clientIpFlowLimit",t)},expression:"form.clientIpFlowLimit"}})],1),r("el-form-item",{attrs:{label:"服务端限流"}},[r("el-input",{attrs:{placeholder:"0表示无限制"},model:{value:e.form.serviceIpFlowLimit,callback:function(t){e.$set(e.form,"serviceIpFlowLimit",t)},expression:"form.serviceIpFlowLimit"}})],1),r("el-form-item",{attrs:{label:"轮询方式"}},[r("el-radio-group",{model:{value:e.form.roundType,callback:function(t){e.$set(e.form,"roundType",t)},expression:"form.roundType"}},[r("el-radio",{attrs:{label:0}},[e._v("random")]),r("el-radio",{attrs:{label:1}},[e._v("round-robin")]),r("el-radio",{attrs:{label:2}},[e._v("weight_round_robin")]),r("el-radio",{attrs:{label:3}},[e._v("ip_hash")])],1)],1),r("el-form-item",{staticClass:"is-required",attrs:{label:"IP列表"}},[r("el-input",{attrs:{type:"textarea",autosize:"",placeholder:"格式：127.0.0.1:80 多条换行"},model:{value:e.form.ipList,callback:function(t){e.$set(e.form,"ipList",t)},expression:"form.ipList"}})],1),r("el-form-item",{staticClass:"is-required",attrs:{label:"权重列表"}},[r("el-input",{attrs:{type:"textarea",autosize:"",placeholder:"格式：50 多条换行"},model:{value:e.form.weightList,callback:function(t){e.$set(e.form,"weightList",t)},expression:"form.weightList"}})],1),r("el-form-item",{attrs:{label:"建立连接超时"}},[r("el-input",{attrs:{placeholder:"单位s，0表示无限制"},model:{value:e.form.upstreamConnectTimeout,callback:function(t){e.$set(e.form,"upstreamConnectTimeout",t)},expression:"form.upstreamConnectTimeout"}})],1),r("el-form-item",{attrs:{label:"获取header超时"}},[r("el-input",{attrs:{placeholder:"单位s，0表示无限制"},model:{value:e.form.upstreamHeaderTimeout,callback:function(t){e.$set(e.form,"upstreamHeaderTimeout",t)},expression:"form.upstreamHeaderTimeout"}})],1),r("el-form-item",{attrs:{label:"链接最大空闲时间"}},[r("el-input",{attrs:{placeholder:"单位s，0表示无限制"},model:{value:e.form.upstreamIdleTimeout,callback:function(t){e.$set(e.form,"upstreamIdleTimeout",t)},expression:"form.upstreamIdleTimeout"}})],1),r("el-form-item",{attrs:{label:"最大空闲链接数"}},[r("el-input",{attrs:{placeholder:"0表示无限制"},model:{value:e.form.upstreamMaxIdle,callback:function(t){e.$set(e.form,"upstreamMaxIdle",t)},expression:"form.upstreamMaxIdle"}})],1),r("el-form-item",[r("el-button",{attrs:{type:"primary",disabled:e.submitButDisabled},on:{click:e.handleSubmit}},[e._v("立即提交")])],1)],1)],1)])],1)],1)},i=[],s=(r("ac1f"),r("5319"),r("a9e3"),r("22ce")),l={name:"ServiceCreateHttp",data:function(){return{isEdit:!1,submitButDisabled:!1,form:{serviceName:"",serviceDesc:"",ruleType:0,rule:"",needHttps:0,needWebsocket:0,needStripUrl:0,urlRewrite:"",headerTransfor:"",roundType:2,ipList:"",weightList:"",upstreamConnectTimeout:"",upstreamHeaderTimeout:"",upstreamIdleTimeout:"",upstreamMaxIdle:"",openAuth:0,blackList:"",whiteList:"",clientIpFlowLimit:"",serviceIpFlowLimit:""}}},created:function(){var e=this.$route.params&&this.$route.params.id;e>0&&(this.isEdit=!0,this.fetchData(e))},methods:{fetchData:function(e){var t=this,r={id:e};Object(s["e"])(r).then((function(e){t.form.id=e.data.info.id,t.form.load_type=e.data.info.load_type,t.form.serviceName=e.data.info.service_name,t.form.serviceDesc=e.data.info.service_desc,t.form.ruleType=e.data.http.rule_type,t.form.rule=e.data.http.rule,t.form.needHttps=e.data.http.need_https,t.form.needWebsocket=e.data.http.need_websocket,t.form.needStripUrl=e.data.http.need_strip_uri,t.form.urlRewrite=e.data.http.url_rewrite.replace(/,/g,"\n"),t.form.headerTransfor=e.data.http.header_transfor.replace(/,/g,"\n"),t.form.roundType=e.data.loadBalance.round_type,t.form.ipList=e.data.loadBalance.ip_list.replace(/,/g,"\n"),t.form.weightList=e.data.loadBalance.weight_list.replace(/,/g,"\n"),t.form.upstreamConnectTimeout=e.data.loadBalance.upstream_connect_timeout,t.form.upstreamHeaderTimeout=e.data.loadBalance.upstream_header_timeout,t.form.upstreamIdleTimeout=e.data.loadBalance.upstream_idle_timeout,t.form.upstreamMaxIdle=e.data.loadBalance.upstream_max_idle,t.form.openAuth=e.data.accessControl.open_auth,t.form.blackList=e.data.accessControl.black_list.replace(/,/g,"\n"),t.form.whiteList=e.data.accessControl.white_list.replace(/,/g,"\n"),t.form.clientIpFlowLimit=e.data.accessControl.clientip_flow_limit,t.form.serviceIpFlowLimit=e.data.accessControl.service_flow_limit})).catch((function(){}))},handleSubmit:function(){var e=this;this.submitButDisabled=!0;var t=Object.assign({},this.form);t.urlRewrite=t.urlRewrite.replace(/\n/g,","),t.headerTransfor=t.headerTransfor.replace(/\n/g,","),t.ipList=t.ipList.replace(/\n/g,","),t.weightList=t.weightList.replace(/\n/g,","),t.blackList=t.blackList.replace(/\n/g,","),t.whiteList=t.whiteList.replace(/\n/g,","),t.serviceIpFlowLimit=Number(t.serviceIpFlowLimit),t.needStripUrl=Number(t.needStripUrl),t.needHttps=Number(t.needHttps),t.needWebsocket=Number(t.needWebsocket),t.clientIpFlowLimit=Number(t.clientIpFlowLimit),t.upstreamConnectTimeout=Number(t.upstreamConnectTimeout),t.upstreamHeaderTimeout=Number(t.upstreamHeaderTimeout),t.upstreamIdleTimeout=Number(t.upstreamIdleTimeout),t.upstreamMaxIdle=Number(t.upstreamMaxIdle),this.isEdit?Object(s["i"])(t).then((function(t){e.submitButDisabled=!1,e.$notify({title:"Success",message:"修改成功",type:"success",duration:2e3})})).catch((function(){e.submitButDisabled=!1})):Object(s["b"])(t).then((function(t){e.submitButDisabled=!1,e.$notify({title:"Success",message:"添加成功",type:"success",duration:2e3})})).catch((function(){e.submitButDisabled=!1}))}}},o=l,n=(r("f304"),r("2877")),c=Object(n["a"])(o,a,i,!1,null,"1f330aac",null);t["default"]=c.exports},dc0b:function(e,t,r){},f304:function(e,t,r){"use strict";r("dc0b")}}]);