(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2a30e9d6"],{a55b:function(t,s,a){"use strict";a.r(s);var e=function(){var t=this,s=t.$createElement,e=t._self._c||s;return e("div",{staticClass:"app flex-row align-items-center bg-dark"},[e("vue-particles",{staticClass:"position-absolute w-100 h-100",attrs:{linesColor:"#dedede",color:"#dedede"}}),e("div",{staticClass:"container"},[e("b-row",{staticClass:"justify-content-center"},[e("b-col",{attrs:{md:"8"}},[e("b-card-group",[e("b-card",{staticClass:"p-4",attrs:{"no-body":""}},[e("b-card-body",[e("b-form",[e("h1",{staticClass:"text-dark"},[t._v("Login")]),e("p",{staticClass:"text-muted"},[t._v("Sign In to your account")]),e("b-input-group",{staticClass:"mb-3"},[e("b-input-group-prepend",[e("b-input-group-text",[e("i",{staticClass:"icon-user"})])],1),e("b-form-input",{staticClass:"form-control",attrs:{type:"text",placeholder:"Username"},model:{value:t.name,callback:function(s){t.name=s},expression:"name"}})],1),e("b-input-group",{staticClass:"mb-4"},[e("b-input-group-prepend",[e("b-input-group-text",[e("i",{staticClass:"icon-lock"})])],1),e("b-form-input",{staticClass:"form-control",attrs:{type:"password",placeholder:"Password"},model:{value:t.passwd,callback:function(s){t.passwd=s},expression:"passwd"}})],1),e("b-row",[e("b-col",{attrs:{cols:"6"}},[e("b-button",{staticClass:"px-4 btn-block",attrs:{variant:"primary"},on:{click:t.login}},[t._v("Login\n                                        ")])],1),e("b-col",{staticClass:"text-right",attrs:{cols:"6"}})],1)],1)],1)],1),e("b-card",{staticClass:"text-white bg-primary py-3 d-md-down-none",staticStyle:{width:"44%"}},[e("b-card-body",{staticClass:"text-center"},[e("div",[e("img",{staticStyle:{width:"3rem"},attrs:{src:a("dd88")}}),e("h4",[t._v("Sudis")]),e("p",[t._v("Distributed supervisor process control system .")]),e("a",{staticClass:"btn btn-primary active mt-3",attrs:{href:"https://github.com/ihaiker/sudis",target:"_blank"}},[t._v("Read More!")])])])],1)],1)],1)],1)],1)],1)},o=[],n=(a("cc57"),{name:"Login",data:function(){return{name:"",passwd:""}},mounted:function(){var t=localStorage.getItem("token");t&&this.$router.push("/admin/dashboard")},methods:{login:function(){var t=this;t.$axios.post("/login",{name:t.name,passwd:t.passwd}).then((function(s){var a=s.token;localStorage.setItem("token",a),t.$toast.success("登录成功！"),t.$router.push("/admin/dashboard")})).catch((function(s){t.$alert(s.message)}))}}}),r=n,i=a("9ca4"),c=Object(i["a"])(r,e,o,!1,null,null,null);s["default"]=c.exports},cc57:function(t,s,a){var e=a("064e").f,o=Function.prototype,n=/^\s*function ([^ (]*)/,r="name";r in o||a("149f")&&e(o,r,{configurable:!0,get:function(){try{return(""+this).match(n)[1]}catch(t){return""}}})},dd88:function(t,s,a){t.exports=a.p+"static/img/logo2.66d14f79.png"}}]);
//# sourceMappingURL=chunk-2a30e9d6.88f246a6.js.map