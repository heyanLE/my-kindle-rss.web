const LeftLoad = $("#left_li_load");
const RightLoad = $("#right_li_load");

let FeedList = [];
let FeedMapKeyClass = {"全部订阅":[],"我的订阅":[]};
let ClassNameList = [];
let UserFeedList = [];

let ClassLiMap = {};
let FromDivMap = {};

let FromNameMapKeyClass = {};

let NowChoseClass = "全部订阅";
let NowChoseFrom = "";

let NowShowFeedList = [];
let NowPage = 1;
let Page = 1;

getFeedInter();

function getFeedInter() {
    FetchGet("/api/v1/feed-list?_type=1",(code,message,value)=>{
        if (code === 200){
            FeedList = value;
            getFeedData();
        } else{
            getFeedInter();
        }
    })
}/*√*/
function getFeedData() {
    let m = FeedList.length;
    newFromSpan("全部");
    for(let i = 0 ; i < m ; i ++){
        const feed = FeedList[i];
        if (FeedMapKeyClass[feed.class] === undefined){
            newClassLi(feed.class);
            FeedMapKeyClass[feed.class] = [feed];
        } else{
            FeedMapKeyClass[feed.class][FeedMapKeyClass[feed.class].length] = feed;
        }
        if (ClassNameList.indexOf(feed.class) === -1){
            ClassNameList[ClassNameList.length] = feed.class;
        }

        if (FromNameMapKeyClass[feed.class] === undefined) {
            FromNameMapKeyClass[feed.class] = ["全部",feed.from];
        }else{
            FromNameMapKeyClass[feed.class][FromNameMapKeyClass[feed.class].length] = feed.from;
        }

        FeedMapKeyClass["全部订阅"][FeedMapKeyClass["全部订阅"].length] = feed;

        if (FromNameMapKeyClass["全部订阅"] === undefined){
            FromNameMapKeyClass["全部订阅"] = ["全部",feed.from];
            newFromSpan(feed.from);
        } else if(FromNameMapKeyClass["全部订阅"].indexOf(feed.from) === -1){
            FromNameMapKeyClass["全部订阅"][FromNameMapKeyClass["全部订阅"].length] = feed.from;
            newFromSpan(feed.from);
        }
    }
    ClassLiMap["全部订阅"] = $("#feed_li_all");
    ClassLiMap["我的订阅"] = $("#feed_li_mine");
    FromNameMapKeyClass["我的订阅"] = ["全部"];
    finishInterLoad();
}/*√*/
function finishInterLoad() {
    dismissLeftLoad();
    onClassLiClick("全部订阅");
}

/*登陆成功调用*/
function getUserFeed() {
    FetchGet("/api/v1/feed",(code , message , value)=>{
        if (code === 200){
            //console.log(value);
            UserFeedList = value;
            onFromSpanClick(NowChoseFrom);
        }
    })
}

function onClassLiClick(id) {
    disAct();
    NowChoseClass = id;
    NowChoseFrom = "";
    act();
    refreshFromDiv();
    refreshFeedList();
}
function onFromSpanClick(id) {
    disAct();
    NowChoseFrom = id;
    if (id === "全部"){
        NowChoseFrom = "";
    }
    act();
    refreshFeedList();
}

function newClassLi(classN) {
    const li = $("<li id=\"feed_li_class_"+classN+"\" class=\"mdui-list-item mdui-ripple\">\n" +
        "  <div class=\"mdui-list-item-content\">"+classN+"</div>\n" +
        "</li>");
    li.click(function () {
        onClassLiClick(classN);
    });
    $("#feed_list_ul").append(li);
    ClassLiMap[classN] = li;
}/*√*/
function newFromSpan(from) {
    const  fromS = $("<button type=\"button\" class=\"mdui-btn\">"+from+"</button>");
    fromS.click(function () {
        onFromSpanClick(from);
    });
    $("#feed_from_group").append(fromS);
    FromDivMap[from] = fromS;
}/*√*/

function refreshFeedList() {
    showRightLoad();

    if (NowChoseClass === "我的订阅"){
        NowShowFeedList = UserFeedList;
        NowPage = 1;
        Page = Math.ceil(NowShowFeedList.length/16) === 0 ? 1 : Math.ceil(NowShowFeedList.length/16);
        refreshFeedListReal();
        return;
    }

    const FeedL = FeedMapKeyClass[NowChoseClass];
    if (FeedL === undefined){
        return;
    }
    const m = FeedL.length;
    NowShowFeedList = [];

    for (let i = 0 ; i < m ; i ++){
        if ((NowChoseFrom !== "" && NowChoseFrom === FeedL[i].from) || NowChoseFrom === ""){
            NowShowFeedList[NowShowFeedList.length] = FeedL[i];
        }
    }

    NowPage = 1;
    Page = Math.ceil(NowShowFeedList.length/16) === 0 ? 1 : Math.ceil(NowShowFeedList.length/16);


    refreshFeedListReal();
}

function refreshFeedListReal() {
    const m = NowShowFeedList.length;
    $("#feed_content_ul").empty();

    const min = 16*(NowPage-1);
    let max = 16*NowPage-1;
    if (max > m-1){
        max = m-1;
    }

    for (let i = min ; i <= max ; i ++){
        newFeedLi(NowShowFeedList[i]);
    }
    refreshPageView();
    dismissRightLoad();
}

function newFeedLi(feed) {
    let s;
    let c;
    if (isSub(feed)) {
        s = "已订阅";
        c = "mdui-color-blue-accent";
    }else{
        s = "订阅";
        c = "mdui-color-pink-accent"
    }
    const li = $("<ul id=\"feed_content_ul\" class=\"mdui-list\">\n" +
        "                            <li class=\"mdui-list-item\">\n" +
        "                                <i class=\"mdui-list-item-avatar mdui-icon material-icons\">&#xe865;</i>\n" +
        "                                <div class=\"mdui-list-item-content\">\n" +
        "                                    <div class=\"mdui-list-item-title\">"+feed.name+"</div>\n" +
        "                                    <div class=\"mdui-list-item-text\" style='margin-right: 40px'>"+feed.describe+"</div>\n" +
        "                                </div>\n" +
        "                                <div id='feed_list_li_button_"+feed.id+"' class=\""+c+" mdui-btn mdui-btn-dense mdui-ripple\" style=\"min-width: 50px;padding: 0\">"+s+"</div>\n" +
        "                            </li>\n" +
        "                        </ul>");
    
    $("#feed_content_ul").append(li);
    $("#feed_list_li_button_"+feed.id).click(function (e) {
        e.stopPropagation();
        onFeedLiButtonClick(feed);
    });
    li.click(function () {
        onFeedLiClick(feed);
    })
}/*√*/

function onFeedLiClick(feed) {
    FetchGet("/api/v1/article?_feed_id="+feed.id,(code,message,value) => {
        //console.log(value[0].Link);
        if (code === 200){
            let aa = "";
            for (let i = 0 ; i < (value.length<9?value.length:9) ; i ++){
                let date = new Date(value[i].PubTime);
                let yy = date.getFullYear();      //年
                let mm = date.getMonth() + 1;     //月
                let dd = date.getDate();          //日
                //console.log(i);
                aa += "<a target='_blank' href=\""+value[i].Link+"\" class=\"mdui-list-item mdui-ripple\">"+yy+"-"+mm+"-"+dd+"  |  "+value[i].Headline+"</a>";
            }
            let bb = "";
            if (!isSub(feed)){
                bb = "<button style='text-align: center !important;' id='dialog_feed_button' class=\"mdui-btn mdui-ripple mdui-text-color-pink-accent\">订阅</button>";
            }else{
                bb = "<button style='text-align: center !important;' id='dialog_feed_button' class=\"mdui-btn mdui-ripple mdui-text-color-blue-accent\">已订阅</button>";
            }

            let diaDia = $("#dialog_feed_dia");
            let dia = $(
                "    <div class=\"mdui-dialog-title\">"+feed.name+"</div>\n" +
                "    <div class=\"mdui-dialog-content\">\n" +
                "        <div class=\"mdui-list\">" + aa + "</div>\n" +
                "    </div>\n" +
                "    <div class=\"mdui-dialog-actions mdui-dialog-actions-stacked\">"+bb+"" +
                "    <button style='text-align: center !important;' id='dialog_feed_button' class=\"mdui-btn mdui-ripple mdui-text-color-pink-accent\" mdui-dialog-close>取消</button>   " +
                "</div>");

            diaDia.empty();
            diaDia.append(dia);

            $("#dialog_feed_button").click(function () {
                onFeedLiButtonClickWC(feed,(b) => {
                    let l = $("#dialog_feed_button");
                    if (b){
                        l.text("已订阅");
                        l.removeClass("mdui-text-color-pink-accent");
                        l.addClass(" mdui-text-color-blue-accent");
                    }else{
                        l.text("订阅");
                        l.removeClass("mdui-text-color-blue-accent");
                        l.addClass(" mdui-text-color-pink-accent");
                    }
                })
            });

            let inst = new mdui.Dialog("#dialog_feed_dia");
            inst.open();
        }
    })
}
function onFeedLiButtonClick(feed) {
    //console.log("ButtonClick");
    if (isSub(feed)){
        //取消订阅
        let ii = 0;
        for(;ii < UserFeedList.length ; ii ++){
            if (UserFeedList[ii].id === feed.id){
                //console.log(ii);
                break;
            }
        }
        FetchDeleteWithBody("/api/v1/feed",{
            "feed-id":feed.id
        },(code,message,value) => {
            if (code === 200){
                UserFeedList.splice(ii,1);
                //console.log(UserFeedList);
                onFromSpanClick(NowChoseFrom);
                mdui.snackbar({
                    message: '已取消订阅',
                    position:`top`
                });
            }else{
                mdui.snackbar({
                    message: message,
                    position:`top`
                });
            }
        })
    }else{
        //订阅
        FetchPost("/api/v1/feed",{
            "feed-id":feed.id
        },(code,message,value) => {
            if (code === 200){
                UserFeedList[UserFeedList.length] = feed;
                onFromSpanClick(NowChoseFrom);
                mdui.snackbar({
                    message: '订阅成功',
                    position:`top`
                });
            }else{
                mdui.snackbar({
                    message: message,
                    position:`top`
                });
            }
        })
    }
    return false;
}
function onFeedLiButtonClickWC(feed,callback) {
    //console.log("ButtonClick");
    if (isSub(feed)){
        //取消订阅
        let ii = 0;
        for(;ii < UserFeedList.length ; ii ++){
            if (UserFeedList[ii].id === feed.id){
                //console.log(ii);
                break;
            }
        }
        FetchDeleteWithBody("/api/v1/feed",{
            "feed-id":feed.id
        },(code,message,value) => {
            if (code === 200){
                UserFeedList.splice(ii,1);
                //console.log(UserFeedList);
                onFromSpanClick(NowChoseFrom);
                mdui.snackbar({
                    message: '已取消订阅',
                    position:`top`
                });
                callback(false);
            }else{
                mdui.snackbar({
                    message: message,
                    position:`top`
                });
            }
        })
    }else{
        //订阅
        FetchPost("/api/v1/feed",{
            "feed-id":feed.id
        },(code,message,value) => {
            if (code === 200){
                UserFeedList[UserFeedList.length] = feed;
                onFromSpanClick(NowChoseFrom);
                mdui.snackbar({
                    message: '订阅成功',
                    position:`top`
                });
                callback(true);
            }else{
                mdui.snackbar({
                    message: message,
                    position:`top`
                });
            }
        })
    }
    return false;
}

function isSub(feed) {
    for (let i = 0 ; i < UserFeedList.length ; i ++){
        if (feed.id === UserFeedList[i].id){
            //console.log("true");
            return true;
        }
    }
    //console.log(feed);
    return false;
}/*√*/

function refreshFromDiv() {
    for (let fromN in FromDivMap){
        if (!FromDivMap.hasOwnProperty(fromN)) {
            continue;
        }
        if (FromNameMapKeyClass[NowChoseClass].indexOf(fromN) === -1){
            FromDivMap[fromN].css("cssText","display:none!important");
        } else{
            FromDivMap[fromN].css("cssText","display:block!important");
        }
    }

}/*√*/

function disAct() {
    if( ClassLiMap[NowChoseClass] !== undefined){
        ClassLiMap[NowChoseClass].removeClass("mdui-list-item-active");
    }
    if (NowChoseFrom === ""){
        FromDivMap["全部"].removeClass("mdui-btn-active");
    }
    if (FromDivMap[NowChoseFrom] !== undefined){
        FromDivMap[NowChoseFrom].removeClass("mdui-btn-active");
    }
}/*√*/

function act() {
    if( ClassLiMap[NowChoseClass] !== undefined){
        ClassLiMap[NowChoseClass].addClass("mdui-list-item-active");
    }
    if (NowChoseFrom === ""){
        FromDivMap["全部"].addClass("mdui-btn-active");
    }
    if (FromDivMap[NowChoseFrom] !== undefined){
        FromDivMap[NowChoseFrom].addClass("mdui-btn-active");
    }
}/*√*/


function showLeftLoad() {
    LeftLoad.css("cssText","display:block!important;width: 60px;height: 60px;margin: 30px")
}/*√*/
function dismissLeftLoad() {
    LeftLoad.css("cssText","display:none!important")
}/*√*/
function showRightLoad() {
    RightLoad.css("cssText","display:block!important;width: 60px;height: 60px;margin: 30px")
}/*√*/
function dismissRightLoad() {
    RightLoad.css("cssText","display:none!important")
}/*√*/


function refreshPageView() {

    clearPageSpan();

    if (NowPage === 1){
        lastButton(false);
    } else{
        lastButton(true);
    }
    if (Page === NowPage) {
        nextButton(false);
    }else{
        nextButton(true);
    }
    
    if (Page <= 9){
        for (let i = 1 ; i <= Page ; i ++){
            addPageButton(i);
        }
    } else{
        if (NowPage <= 5){
            addPageButton(1);
            addPageButton(2);
            addPageButton(3);
            addPageButton(4);
            addPageButton(5);
        } else{
            addPageButton(1);
            addMSpan();

            let n = NowPage;
            if (n >= Page - 4) {
                n = Page - 4;
            }

            addPageButton(n - 2);
            addPageButton(n - 1);
            addPageButton(n);
        }

        if (NowPage >= Page - 4){
            addPageButton(Page - 3);
            addPageButton(Page - 2);
            addPageButton(Page - 1);
            addPageButton(Page);
        } else{
            if (NowPage <= 5){
                addPageButton(6);
                addPageButton(7);
                addMSpan();
                addPageButton(Page);
            } else{
                addPageButton(NowPage + 1);
                addPageButton(NowPage + 2);
                addMSpan();
                addPageButton(Page);
            }
        }
    }
    
}/*√*/

function nextButton(isShow) {
    if (isShow){
        $("#button_page_next").css("display","inline-block");
    } else{
        $("#button_page_next").css("display","none");
    }
}/*√*/
function lastButton(isShow) {
    if (isShow){
        $("#button_page_last").css("display","inline-block");
    } else{
        $("#button_page_last").css("display","none");
    }
}/*√*/

function onPageNextButtonClick() {
    if (NowPage < Page) {
        onPageButtonClick(NowPage + 1);
    }
}/*√*/
function onPageLastButtonClick() {
    if (NowPage > 1) {
        onPageButtonClick(NowPage - 1);
    }
}/*√*/

function clearPageSpan() {
    $("#page_button_span").empty();
}/*√*/

function addPageButton(page) {
    const b = $("<button class=\"mdui-btn mdui-btn-dense mdui-ripple\" style=\"text-align: center;min-width: 32px;width: 32px;height: 32px;margin: 0;padding: 0\">"+page+"</button>");
    $("#page_button_span").append(b);
    b.click(function () {
        onPageButtonClick(page);
    });
    if (page === NowPage){
        b.addClass("btn-active");
    }
}/*√*/

function addMSpan() {
    const m = $("<span> ... </span>");
    $("#page_button_span").append(m);
}/*√*/

function onPageButtonClick(page) {
    NowPage = page;
    refreshFeedListReal();
    backTop();
}

function backTop() {
    //console.log("Top");
    (function s(){
        let currentScroll = document.documentElement.scrollTop || document.body.scrollTop;
        if (currentScroll > 0) {
            window.requestAnimationFrame(s);
            window.scrollTo (0,currentScroll - (currentScroll/5));
        }
    })();

    window.scrollTo(0,0);

    $("#context").scrollTop(0);
}
