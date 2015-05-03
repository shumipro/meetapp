import $ from 'jquery'
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

        render()
        // TODO: check logged in

        // move to login

        // attach post

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

    render() {
        this._$wrap.empty()
        for(var i = 0; i < list.length; i++) {
            var $item = $(tmpl({
                id: list[i].id,
                name: list[i].name,
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

    }
}