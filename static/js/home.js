let User = null;

let LoginCaptchaId = 0;
let RegisterCaptchaId = 0;

let RegisterDataUnix = 1563367986;
let RegisterToken = "";
let RegisterEmail = "";
let RegisterEmailCaptcha = "";

let RegisterEmailSendTimer; //timer变量，控制时间
let RegisterEmailSendCount = 60; //间隔函数，1秒执行
let RegisterEmailSendCurCount;//当前剩余秒数

let LoginDialog = null;
let RegisterDialogFirst = null;
let RegisterDialogNext = null;

FetchGet("/api/v1/user",(code,message,value)=>{
    if (code === 200){
        User = value;
        refreshUser();
    }
});

/*回车事件*/
$('#login_email_input').bind('keypress', function(event) {
    console.log("KeyCode"+event.keyCode);
    if (event.keyCode == "13") {
        console.log("KeyCode"+event.keyCode);
        event.preventDefault();
        //回车执行查询
        $('#login_button').click();
        return false
    }
});
$('#login_password_input').bind('keypress', function(event) {
    if (event.keyCode == "13") {
        event.preventDefault();
        //回车执行查询
        $('#login_button').click();
        return false
    }
});
$('#login_captcha_input').bind('keypress', function(event) {
    if (event.keyCode == "13") {
        event.preventDefault();
        //回车执行查询
        $('#login_button').click();
        return false
    }
});

function onRegisterPasswordInputOnInput(){
    console.log($("#register_password_input").val().length);
    if ($("#register_password_input").val().length <= 6){
        $("#register_next_password_div").addClass("mdui-textfield-invalid");
    } else{
        $("#register_next_password_div").removeClass("mdui-textfield-invalid");
    }
}

function onRegisterPasswordInputCInput(){
    if ($("#register_password_input").val() !== $("#register_password_c_input").val()){
        $("#dialog_register_password_c_div").addClass("mdui-textfield-invalid");
    } else{
        $("#dialog_register_password_c_div").removeClass("mdui-textfield-invalid");
    }
}



/*=====================Utils==================*/
function NewLoginCaptcha() {
    FetchGet("/api/v1/captcha",(code,message,value) => {
        LoginCaptchaId = value.captcha_id;
        const img = $("#dialog_login_captcha_img");
        $("#dialog_login_captcha_div").css("display" ,"block");
        img.attr("src","/api/v1/captcha/img/"+LoginCaptchaId+".png");
        $("#dialog_login").css("overflow-y","auto");
        img.bind("click",function(){
            console.log("ImgClick");
            onLoginCaptchaImgClick()
        });

    })
}/*√*/

function NewRegisterCaptcha() {
    FetchGet("/api/v1/captcha",(code,message,value) => {
        RegisterCaptchaId = value.captcha_id;
        $("#dialog_register_captcha_D_div").css("display" ,"block");
        $("#dialog_register_captcha_img").attr("src","/api/v1/captcha/img/"+RegisterCaptchaId+".png");
        $("#dialog_register_first").css("overflow-y","auto");
        $("#dialog_register_captcha_img").bind("click",function(){
            console.log("ImgClick");
            onRegisterCaptchaImgClick()
        });
    })
}/*√*/

/**
 * @return {string}
 */
function GetNow() {
    return  (Date.parse( new Date() )/1000).toString();
}

function refreshUser() {
    if (User == null){
        $("#user_card_content").text("未登录");
        $("#balance_card_content").text(0);
        $("#push_card_content").text("自动推送关闭");
        $("#user_card_action_out").css("display","block");
        $("#user_card_action_up").css("display","none");
    } else{
        $("#user_card_content").text(User.email);
        $("#balance_card_content").text(User.balance);
        if (User.push_auto){
            $("#push_card_content").text("自动推送:"+User.push_time+"点");
        } else{
            $("#push_card_content").text("自动推送关闭");
        }
        $("#user_card_action_out").css("display","none");
        $("#user_card_action_up").css("display","block");
    }
    const w = $("#wallpaper");
    w.fadeToggle(1);
    w.fadeToggle(800);
}/*√*/
/*=================Click======================*/

/*Content*/
function onContentLoginButtonClick() {
    if(LoginCaptchaId === 0){
        $("#dialog_login_captcha_div").css("display" ,"none");
        $("#dialog_login_captcha_img").attr("src","/api/v1/captcha/img/"+LoginCaptchaId+".png");
        $("#dialog_login").css("overflow-y","none");
        LoginDialog = new mdui.Dialog("#dialog_login", "models:true");
        LoginDialog.open();
    }else{
        NewLoginCaptcha();
        LoginDialog = new mdui.Dialog("#dialog_login", "models:true");
        LoginDialog.open();
    }
}/*√*/

function onContentRegisterButtonClick() {
    if(RegisterCaptchaId === 0){
        $("#dialog_register_captcha_D_div").css("display" ,"none");
        $("#dialog_register_captcha_img").attr("src","/api/v1/captcha/img/"+LoginCaptchaId+".png");
        $("#dialog_register_first").css("overflow-y","none");
        RegisterDialogFirst = new mdui.Dialog("#dialog_register_first", "models:true");
        RegisterDialogFirst.open();
    }else{
        NewLoginCaptcha();
        RegisterDialogFirst = new mdui.Dialog("#dialog_register_first", "models:true");
        RegisterDialogFirst.open();
    }
}/*√*/

function onContentLoginOutButtonClick() {
    FetchDelete("/api/v1/user",()=>{
        User = null;
        refreshUser();
    });
    mdui.snackbar({
        message: '登出成功',
        position:`top`
    });
}/*√*/

/*LoginDialog*/
function onLoginPasswordInputClick() {
    $("#dialog_login_password_error").text("密码不能为空");
    $("#dialog_login_password_div").removeClass("mdui-textfield-invalid");
}/*√*/

function onLoginCaptchaImgClick() {
    $("#dialog_login_captcha_img").attr("src", "/api/v1/captcha/img/"+LoginCaptchaId+".png?reload="+GetNow);
}/*√*/

function onLoginCaptchaInputClick() {
    $("#login_captcha_err").text("密码不能为空");
    $("#dialog_login_captcha_input_div").removeClass("mdui-textfield-invalid");
}/*√*/

function onLoginButtonClick() {
    const loginEmailDiv = $("#dialog_login_email_input_div");
    if (loginEmailDiv.hasClass("mdui-textfield-invalid")
    || loginEmailDiv.hasClass("mdui-textfield-invalid-html5")){
        return;
    }
    const email = $("#login_email_input").val();
    const password = $("#login_password_input").val();
    if (email === ""){
        loginEmailDiv.addClass("mdui-textfield-invalid");
        return;
    }
    if (password === ""){
        $("#dialog_login_password_div").addClass("mdui-textfield-invalid");
        return;
    }
    if(LoginCaptchaId === 0){
        FetchPost("/api/v1/user",{
            email:email,
            password:password
        },(code,message,value)=>{
            switch (code) {
                case 200:
                    LoginDialog.close();
                    User = value;
                    refreshUser();
                    mdui.snackbar({
                        message: '登陆成功',
                        position:`top`
                    });
                    break;
                case 412:
                    NewLoginCaptcha();
                    LoginDialog.handleUpdate();
                    break;
                default:
                    $("#dialog_login_password_error").text(message);
                    $("#dialog_login_password_div").addClass("mdui-textfield-invalid");
            }
        })
    }else{
        const c = $("#login_captcha_input");
        const captcha = c.val();
        if (captcha === ""){
            c.addClass("mdui-textfield-invalid");
            return;
        }
        FetchPostWithCaptcha("/api/v1/user",{
            email:email,
            password:password
        },LoginCaptchaId,captcha,(code,message,value)=>{
            switch (code) {
                case 200:
                    LoginDialog.close();
                    User = value;
                    refreshUser();
                    mdui.snackbar({
                        message: '登陆成功',
                        position:`top`
                    });
                    break;
                case 412:
                    $("#login_captcha_err").text("验证码错误");
                    $("#dialog_login_captcha_input_div").addClass("mdui-textfield-invalid");
                    NewLoginCaptcha();
                    LoginDialog.handleUpdate();
                    break;
                default:
                    $("#dialog_login_password_error").text(message);
                    $("#dialog_login_password_div").addClass("mdui-textfield-invalid");
            }
        })
    }
}/*√*/

/*RegisterFirstDialog*/

function onRegisterFirstEmailInputClick() {
    $("#register_email_err").text("请输入正确的邮箱");
    $("#dialog_register_first_email_div").removeClass("mdui-textfield-invalid");
}/*√*/

function onRegisterEmailCaptchaInputClick() {
    $("#register_email_captcha_err").text("验证码不能为空");
    $("#register_email_captcha_div").removeClass("mdui-textfield-invalid");
}/*√*/

function onRegisterSendEmailButtonClick() {
    const registerEmailDiv = $("#dialog_register_first_email_div");
    if (registerEmailDiv.hasClass("mdui-textfield-invalid")
        || registerEmailDiv.hasClass("mdui-textfield-invalid-html5")) {
        return;
    }

    const email = $("#register_email_input").val();
    if (email === ""||email.indexOf("@")=== -1||email.indexOf(".") === -1){
        registerEmailDiv.addClass("mdui-textfield-invalid");
    }

    const button = $("#send_email_button");
    button.attr("disabled",true);
    button.text("发送中...");


    let url = "";
    url = "/api/v1/register?_email="+email;
    if (RegisterCaptchaId !== 0) {
        const captcha = $("#register_captcha_input").val();
        if (captcha === ""){
            button.removeAttr("disabled");//启用按钮
            button.val("点击发送");
            return;
        }
        console.log("Register_captcha_id =>"+RegisterCaptchaId);
        console.log("Register_captcha =>"+captcha);
        url += "&_captcha_id="+RegisterCaptchaId+"&_captcha="+captcha;
    }
    console.log("Fetch-register");

    FetchGet(url, (code,message,value) => {
        switch (code) {
            case 200:
                RegisterEmail = value.email;
                RegisterToken = value.token;
                RegisterDataUnix = value.date_unix;
                button.attr("disabled",true);
                RegisterEmailSendCurCount = RegisterEmailSendCount;
                button.text(RegisterEmailSendCurCount+"s");
                RegisterEmailSendTimer = window.setInterval(SetRemainTime, 1000);
            function SetRemainTime() {
                if (RegisterEmailSendCurCount === 0) {
                    window.clearInterval(RegisterEmailSendTimer);//停止计时器
                    const button = $("#send_email_button");
                    button.removeAttr("disabled");//启用按钮
                    button.text("点击发送");
                }
                else {
                    button.attr("disabled",true);//启用按钮
                    RegisterEmailSendCurCount--;
                    $("#send_email_button").text(RegisterEmailSendCurCount+"s");
                }
            }
                break;
            case 412:
                if (RegisterCaptchaId !== 0) {
                    $("#register_captcha_err").text("验证码错误");
                    $("#register_captcha_input_div").addClass("mdui-textfield-invalid")
                }else{
                    mdui.snackbar({
                        message: '本次发送需要输入人工验证码',
                        position:`top`
                    });
                }
                button.removeAttr("disabled");
                button.text("点击发送");
                NewRegisterCaptcha();
                RegisterDialogFirst.handleUpdate();
                break;
            case 231:
                $("#register_email_err").text("邮箱已存在");
                button.removeAttr("disabled");
                button.text("点击发送");
                $("#dialog_register_first_email_div").addClass("mdui-textfield-invalid");
                break;
            default:
                $("#register_captcha_err").text(message);
                $("#register_captcha_input_div").addClass("mdui-textfield-invalid");
                button.removeAttr("disabled");
                button.text("点击发送");
        }
    });
}/*√*/

function onRegisterCaptchaInputClick() {
    $("#register_captcha_err").text("验证码不能为空");
    $("#register_captcha_input_div").removeClass("mdui-textfield-invalid");
}/*√*/

function onRegisterCaptchaImgClick() {
    console.log("Img Click");
    $("#dialog_register_captcha_img").attr("src", "/api/v1/captcha/img/"+RegisterCaptchaId+".png?reload="+GetNow);
}/*√*/

function onRegisterNextButtonClick() {

    const captchaDiv = $("#register_email_captcha_div");
    const captcha = $("#register_email_captcha_input").val();
    if (captcha === ""){
        $("#register_captcha_err").text("验证码不能为空");
        captchaDiv.addClass("mdui-textfield-invalid");
        return;
    }

    if (RegisterEmail === "" || RegisterToken === ""){
        return
    }

    let p = "?_email="+RegisterEmail+"&_token="+RegisterToken+"&_date_unix="+RegisterDataUnix+"&_captcha="+captcha;
    FetchGet("/api/v1/register-verify"+p,(code,message,value)=>{
        if (code === 200){
            RegisterEmailCaptcha = captcha;
            RegisterDialogFirst.close();
            RegisterDialogNext = new mdui.Dialog("#dialog_register_next", "models:true");
            RegisterDialogNext.open();
        } else{
            $("#register_email_captcha_err").text(message);
            captchaDiv.addClass("mdui-textfield-invalid");
        }
    });
}/*√*/



/*RegisterDialog*/
function onRegisterButtonClick() {

    const passwordF = $("#register_password_input").val();
    const passwordN = $("#register_password_c_input").val();

    if (passwordF.length < 6){
        $("#register_next_password_div").addClass("mdui-textfield-invalid");
        return;
    }
    if (passwordN !== passwordF){
        $("#dialog_register_password_c_div").addClass("mdui-textfield-invalid");
        return;
    }

    const date = {
        email:RegisterEmail,
        password:passwordN,
        token:RegisterToken,
        date_unix:RegisterDataUnix,
        captcha:RegisterEmailCaptcha
    };

    FetchPost("/api/v1/register",date,(code,message,value)=>{
        if (code === 200){
            User = value;
            refreshUser();
            mdui.snackbar({
                message: '注册成功，已自动登录',
                position:`top`
            });
            RegisterDialogNext.close();
        } else{
            RegisterDialogNext.close();
            mdui.snackbar({
                message: '注册失败，'+ message,
                position:`top`
            });
        }
    })
}/*√*/
