import $ from 'jquery'
import util from '../util'

export default class Notifications {
    constructor() {
        this._$wrap = $('.ma-nav-notification')
        this._$link = $('.ma-nav-notification-link')
        this._$badge = this._$link.find('.badge')
        this._$dropdown = $('.ma-nav-notification-dropdown')
        // attach events
        this._$link.on('click', ()=> {
            this.read()
        })
        // load data
        if(util.getUserInfo()){
            this.load()
        }
    }

    load(){
        $.ajax({
            url: '/u/api/notification'
        }).done((res) => {
            console.log(res)
            this.render(res)
        }).fail(() => {
            alert("Error")
        })
    }

    read(){
        $.ajax({
            url: '/u/api/notification/done',
            type: 'put',
            contentType:"application/json; charset=utf-8",
            dataType: 'json'
        }).done((res) => {
            this._$badge.removeClass('ma-unread-badge').html(0)
        }).fail(() => {
            alert("Error")
        })
    }

    render(res){
        if(!res && !res.Notifications){
            return
        }
        /**
         * DetailURL: "/app/detail/2e885e1a-01fb-4138-a2aa-6b4350b03ea4"
         * IsRead: false
         * Message: "新着メッセージ: aa"
         * NotificationID: "1431247494803451353"
         * NotificationType: 1
         * SourceID: "1431247494803451353"
         */
        var list = res.Notifications,
            unreadCount = 0
        for(var i=list.length - 1; i>= 0; i--){
            var item = list[i]
            if(!item.IsRead){
                unreadCount++
            }
            var displayMsg = item.Message
            if(displayMsg.length > 20){
                displayMsg = displayMsg.substring(0, 20) + "..."
            }
            var $li = $('<li role="presentation"><a href="' + item.DetailURL + '">' + displayMsg + '</a></li>')
            $li.appendTo(this._$dropdown)
        }
        // unread exists
        if(unreadCount > 0){
            this._$badge.addClass('ma-unread-badge').html(unreadCount)
        }
        this._$wrap.css('visibility', 'visible')
    }
}