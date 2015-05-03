import $ from 'jquery'
import util from '../util'
import Handlebars from 'handlebars'

var html = '<li class="ma-profile-comment">' +
                '<span class="ma-profile-date subtle-text">{{date}}</span>' +
                '<a href="#">' +
                    '<div class="ma-profile-image-wrap">' +
                        '<img class="img-rounded ma-profile-image" src="{{imageUrl}}}">' +
                    '</div>' +
                '</a>' +
                '<div class="ma-profile-message">' +
                    '<p>{{message}}</p>' +
                '</div>' +
            '</li>'
var tmpl = Handlebars.compile(html)

export default class ConstantList {
    constructor() {
        this._$wrap = $('.ma-discussion-wrap')
        this._$textarea = $('#ma_detail_discussion_comment')
        $('.ma_detail_discussion_comment_btn').on('click', () => {
            if(util.getUserInfo()){
                this.post()
            }else{
                // move to login for anonymous
                location.href = "/login"
            }
        })
        // temp
        var discussions = window.ma_data && window.ma_data.discussions
        render(discussions)
    }

    load() {
        // GET /api/app/discussions?appId=0d58a936-f148-404e-91a5-1762331367d8
        // response:
        // [{
        //     appId: "0d58a936-f148-404e-91a5-1762331367d8",
        //     userId: "10152160532855662",
        //     userName: "Takuya Tejima",
        //     message: "hoge hoge",
        //     timestamp: "1430622775846"
        // },
        // {
        //     appId: "0d58a936-f148-404e-91a5-1762331367d8",
        //     userId: "10152160532855662",
        //     userName: "Takuya Tejima",
        //     message: "hoge hoge2222",
        //     timestamp: "1430622775946"
        // }]


    }

    render(discussions) {
        this._$wrap.empty()
        for(var i = 0; i < discussions.length; i++) {
            var $item = $(tmpl({
                id: discussions[i].id,
                name: discussions[i].name,
                prop: prop
            }))
            $item.appendTo(this._$wrap)
        }
    }

    post() {
        // POST /api/app/discussion
        // request body (json):
        // {
        //     appId: "0d58a936-f148-404e-91a5-1762331367d8",
        //     userId: "10152160532855662",
        //     message: "hoge hoge",
        //     timestamp: "1430622775846"
        // }
        // response:
        // {
        //     appId: "0d58a936-f148-404e-91a5-1762331367d8",
        //     userId: "10152160532855662",
        //     userName: "Takuya Tejima",
        //     message: "hoge hoge",
        //     timestamp: "1430622775846"
        // }
        var params = this.getParams()
        //  validation
        var result = this.validate(params)
        if(result.error){
            alert(result.message)
            return
        }
        // {"name": "App name", "description": "hoge", "images": [{"url": "https://golang.org/doc/gopher/gopherbw.png"}]}
        $.ajax({
            url: '/api/app/discussion',
            type: 'post',
            contentType:"application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify(params)
        }).done((res) => {
            // TODO: append result
            alert(JSON.stringify(res))
        }).fail(() => {
            alert("Error")
        })
    }

    validate(params) {
        if($.trim(params.message) === ""){
            return {"error": true, "message": "messageが入力されていません"}
        }
        return {"error": false}
    }

    getParams() {
        return {
            appId: util.getAppDetailId(),
            discussionInfo: {
                userId: util.getUserInfo().id,
                message: this._$textarea.val(),
                timestamp: new Date().getTime()
            }
        }
    }
}