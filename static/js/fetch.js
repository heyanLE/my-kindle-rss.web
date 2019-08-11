function FetchGet(url,callBack) {
    fetch(url,{ method: 'GET'})
        .then((response) => {
            return response.json();
        }).then((responseJson) => {
            callBack(responseJson.code,responseJson.message,responseJson.value)
    })
}
function FetchPost(url,body,callBack) {
    fetch(url,{
        method: `POST`,
        headers: new Headers({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Accept-Language': 'zh-cn,zh;q=0.5',
        }),
        body : JSON.stringify(body)
    }).then((response) => {
        return response.json();
    }).then((responseJson) => {
        callBack(responseJson.code,responseJson.message,responseJson.value)
    })
}
function FetchPostWithCaptcha(url,body,captchaId,captcha,callBack) {
    fetch(url,{
        method: `POST`,
        headers: new Headers({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Accept-Language': 'zh-cn,zh;q=0.5',
            'X-Captcha-Id':captchaId.toString(),
            'X-Captcha':captcha.toString(),
        }),
        body : JSON.stringify(body)
    }).then((response) => {
        return response.json();
    }).then((responseJson) => {
        callBack(responseJson.code,responseJson.message,responseJson.value)
    })
}
function FetchDelete(url,callBack) {
    fetch(url,{
        method: `POST`,
        headers: new Headers({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Accept-Language': 'zh-cn,zh;q=0.5',
            'X-HTTP-Method-Override': 'Delete',
        }),
    }).then((response) => {
        return response.json();
    }).then((responseJson) => {
        callBack(responseJson.code,responseJson.message,responseJson.value)
    })
}

function FetchDeleteWithBody(url,body,callBack) {
    fetch(url,{
        method: `POST`,
        headers: new Headers({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Accept-Language': 'zh-cn,zh;q=0.5',
            'X-HTTP-Method-Override': 'Delete',
        }),
        body : JSON.stringify(body)
    }).then((response) => {
        return response.json();
    }).then((responseJson) => {
        callBack(responseJson.code,responseJson.message,responseJson.value)
    })
}