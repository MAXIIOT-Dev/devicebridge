webpackJsonp([1],{"0U+W":function(t,e){},L0GV:function(t,e){},NHnr:function(t,e,l){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var n=l("7+uW"),i={render:function(){var t=this.$createElement,e=this._self._c||t;return e("div",{attrs:{id:"app"}},[e("router-view")],1)},staticRenderFns:[]};var s=l("VU/8")({name:"App"},i,!1,function(t){l("0U+W")},null,null).exports,a=l("/ocq"),p={name:"HelloWorld",data:function(){return{pldurl:"",pldshow:[{show:!1}],msg:"Welcome to Your Vue.js App",pldlist:[{pp:"device_eui：",ipt:""},{pp:"name：",ipt:""}],pldsblist:[{pp:"设备名称"},{pp:"状态"},{pp:"操作"}],pldfinval:[{pp:"angus-008",status:"online",show:!0,pp1:"编辑",pp2:"删除",dev_eui:"",pldval:""},{pp:"angus-008",status:"online",show:!0,pp1:"编辑",pp2:"删除",dev_eui:"",pldval:""},{pp:"angus-008",status:"online",show:!0,pp1:"编辑",pp2:"删除",dev_eui:"",pldval:""}]}},methods:{pldgeturl:function(){var t=this;t.$ajax.get(t.pld_global.pldbaseurl+"/mappage",{headers:{Authorization:window.localStorage.getItem("token")}}).then(function(e){t.pldurl=e.data.mappage}).catch(function(t){alert(t)})},pldtext:function(t){var e=this;this.pldfinval[t].show?(this.pldfinval[t].show=!1,this.pldfinval[t].pp1="确定"):(this.pldfinval[t].show=!0,this.pldfinval[t].pp1="编辑",e.$ajax.put(e.pld_global.pldbaseurl+"/device",{device_eui:e.pldfinval[t].dev_eui,name:e.pldfinval[t].pldval},{headers:{Authorization:window.localStorage.getItem("token")}}).then(function(){alert("修改成功"),e.getpldsb()}).catch(function(t){alert(t)}))},pldadd:function(){this.pldshow[0].show=!0},pldsure:function(){var t=this;t.$ajax.post(t.pld_global.pldbaseurl+"/device",{device_eui:t.pldlist[0].ipt,name:t.pldlist[1].ipt},{headers:{Authorization:window.localStorage.getItem("token")}}).then(function(){alert("新增成功"),t.pldshow[0].show=!1,t.getpldsb()}).catch(function(t){alert(t)})},suredel:function(t){var e=this;1==confirm("确认要删除吗？")?e.$ajax.delete(e.pld_global.pldbaseurl+"/device/"+e.pldfinval[t].dev_eui,{headers:{Authorization:window.localStorage.getItem("token")}}).then(function(){alert("删除成功"),e.getpldsb()}).catch(function(t){alert(t)}):alert("中途取消")},getpldsb:function(){var t=this;t.$ajax.get(t.pld_global.pldbaseurl+"/device",{params:{page:1,perpage:10},headers:{Authorization:window.localStorage.getItem("token")}}).then(function(e){console.log(e),t.pldfinval.splice(0,t.pldfinval.length);for(var l=0;l<e.data.result.devices.length;l++){var n={pp:"",status:"",show:!0,pp1:"编辑",pp2:"删除",dev_eui:"",pldval:""};n.pp=e.data.result.devices[l].device_name,n.status=e.data.result.devices[l].status,n.dev_eui=e.data.result.devices[l].device_eui,t.pldfinval.push(n)}}).catch(function(){})}},mounted:function(){this.getpldsb(),this.pldgeturl()}},o={render:function(){var t=this,e=t.$createElement,l=t._self._c||e;return l("div",{staticClass:"hello"},[l("iframe",{attrs:{id:"iframe",src:t.pldurl}}),t._v(" "),l("div",{staticClass:"pldright"},[l("img",{staticClass:"pldadmin",attrs:{src:"static/img/admin.png",alt:"admin"}}),t._v(" "),l("div",[l("div",[l("div",{staticClass:"pldaddsb"},[l("button",{on:{click:t.pldadd}},[t._v("添加设备")])]),t._v(" "),l("ul",{directives:[{name:"show",rawName:"v-show",value:t.pldshow[0].show,expression:"pldshow[0].show"}],staticClass:"pldadd"},[t._l(t.pldlist,function(e,n){return l("li",[t._v("\n\t\t\t\t\t\t"+t._s(e.pp)+"\n\t\t\t\t\t\t"),l("input",{directives:[{name:"model",rawName:"v-model",value:e.ipt,expression:"it.ipt"}],attrs:{type:"text",placeholder:"请输入需要的值"},domProps:{value:e.ipt},on:{input:function(l){l.target.composing||t.$set(e,"ipt",l.target.value)}}})])}),t._v(" "),l("button",{on:{click:t.pldsure}},[t._v("确定")])],2),t._v(" "),l("div",{staticClass:"pldlblist"},[l("table",[l("thead",t._l(t.pldsblist,function(e,n){return l("td",[t._v("\n\t\t\t\t\t\t\t\t"+t._s(e.pp)+"\n\t\t\t\t\t\t\t")])}),0),t._v(" "),l("tbody",t._l(t.pldfinval,function(e,n){return l("tr",[e.show?l("td",[t._v("\n\t\t\t\t\t\t\t\t\t"+t._s(e.pp)+"\n\t\t\t\t\t\t\t\t")]):0==e.show?l("td",[l("input",{directives:[{name:"model",rawName:"v-model",value:e.pldval,expression:"ii.pldval"}],attrs:{type:"text",placeholder:"请输入需要的值"},domProps:{value:e.pldval},on:{input:function(l){l.target.composing||t.$set(e,"pldval",l.target.value)}}})]):t._e(),t._v(" "),l("td",[t._v("\n\t\t\t\t\t\t\t\t\t"+t._s(e.status)+"\n\t\t\t\t\t\t\t\t")]),t._v(" "),l("td",[l("button",{on:{click:function(e){return t.pldtext(n)}}},[t._v("\n\t\t\t\t\t\t\t\t\t\t"+t._s(e.pp1)+"\n\t\t\t\t\t\t\t\t\t")]),t._v(" "),l("button",{on:{click:function(e){return t.suredel(n)}}},[t._v("\n\t\t\t\t\t\t\t\t\t\t"+t._s(e.pp2)+"\n\t\t\t\t\t\t\t\t\t")])])])}),0)])])])])])])},staticRenderFns:[]};var d=l("VU/8")(p,o,!1,function(t){l("kl7V")},"data-v-21e95a6f",null).exports,r={name:"login",data:function(){return{ask:!0,userlist:[{pp:""},{pp:""}],pldloginlist:[{pp:"登录"},{pp:"账号或密码",pldkg:!1},{pp:"记住密码"},[{pla:"请输入账号"},{pla:"请输入密码"}]],pldzz:""}},mounted:function(){var t=this;document.onkeydown=function(e){var l=e||window.event||arguments.callee.caller.arguments[0];l&&13==l.keyCode&&t.gettoken()}},methods:{show:function(t){0==this.ask?this.ask=!0:this.ask=!1},gettoken:function(){var t=this;if(""==this.userlist[0].pp||""==this.userlist[1].pp)return this.pldloginlist[1].pp="账号或密码不能为空",void(this.pldloginlist[1].pldkg=!0);this.$ajax.post(this.pld_global.pldbaseurl+"/user/login",{user_name:this.userlist[0].pp,password:this.userlist[1].pp}).then(function(e){window.localStorage.setItem("token",e.data.result.jwt),t.$router.push({path:"/mains"})}).catch(function(e){t.pldloginlist[1].pp="账号或密码错误",t.pldloginlist[1].pldkg=!0})}}},u={render:function(){var t=this,e=t.$createElement,l=t._self._c||e;return l("div",{staticClass:"pldlogin"},[l("div",{staticClass:"pldlogincontent"},[l("div",{staticClass:"pldadddiv"},[l("p",{staticClass:"pldloginp"},[t._v("\n\t\t\t\t"+t._s(t.pldloginlist[0].pp)+"\n\t\t\t")]),t._v(" "),l("div",{staticClass:"pldpld"},[l("p",{directives:[{name:"show",rawName:"v-show",value:t.pldloginlist[1].pldkg,expression:"pldloginlist[1].pldkg"}]},[t._v(t._s(t.pldloginlist[1].pp))]),t._v(" "),t._l(t.pldloginlist[3],function(e,n){return 0==n?l("input",{directives:[{name:"model",rawName:"v-model",value:t.userlist[0].pp,expression:"userlist[0].pp"}],attrs:{type:"text",placeholder:e.pla},domProps:{value:t.userlist[0].pp},on:{input:function(e){e.target.composing||t.$set(t.userlist[0],"pp",e.target.value)}}}):t._e()}),t._v(" "),t._l(t.pldloginlist[3],function(e,n){return 1==n?l("input",{directives:[{name:"model",rawName:"v-model",value:t.userlist[1].pp,expression:"userlist[1].pp"}],attrs:{type:"password",placeholder:e.pla},domProps:{value:t.userlist[1].pp},on:{input:function(e){e.target.composing||t.$set(t.userlist[1],"pp",e.target.value)}}}):t._e()})],2),t._v(" "),l("div",{staticClass:"pldbtn"},[l("button",{on:{click:t.gettoken}},[t._v("\n\t\t\t\t\t"+t._s(t.pldloginlist[0].pp)+"\n\t\t\t\t")])])])]),t._v(" "),l("p",{staticClass:"pldfooterp"},[t._v("\n\t\tCopyright © 2018 ATENXA. ALL Rights Reserved\n\t")])])},staticRenderFns:[]};var c=l("VU/8")(r,u,!1,function(t){l("rlMO")},null,null).exports;n.a.use(a.a);var v=new a.a({routes:[{path:"/mains",name:"HelloWorld",component:d},{path:"/",name:"login",component:c}]}),h=l("aozt"),g=l.n(h),f={pldbaseurl:window.location.protocol+"//"+window.location.host+"/api"},m={render:function(){var t=this.$createElement;return(this._self._c||t)("div")},staticRenderFns:[]};var _=l("VU/8")(f,m,!1,function(t){l("L0GV")},null,null).exports;n.a.config.productionTip=!1,n.a.prototype.$ajax=g.a,n.a.prototype.pld_global=_,new n.a({el:"#app",router:v,components:{App:s},template:"<App/>"})},kl7V:function(t,e){},rlMO:function(t,e){}},["NHnr"]);
//# sourceMappingURL=app.438724879fc5aea11f7f.js.map