import $ from 'jquery'
import util from '../util'

export default class StarButton {
    constructor() {
        // star button
        var $starBtn = $('.ma-app-star-btn')
        $starBtn.on('click', () => {
            // check the user is already logged in
            if(!util.getUserInfo()) {
                // move to login for anonymous
                location.href = '/login'
                return;
            }
            // send star
            $.ajax({
                url: '/u/api/app/star/' + util.getAppDetailId(),
                type: $starBtn.data('api') === 'star' ? 'post' : 'delete',
                contentType:"application/json; charset=utf-8",
                dataType: 'json'
            }).done((res) => {
                location.reload()
            }).fail(() => {
                alert("Error")
            })           
        })
    }
}
