import $ from 'jquery'
import util from '../util'

export default class AppButtons {
    constructor() {
        // common button widgets
        $('.ma-app-edit-btn').on('click', () => {
            location.href = '/u/app/edit/' + util.getAppDetailId()
        })

        $('.ma-app-delete-btn').on('click', () => {
            if(window.confirm('この開発アイデアを削除してもよろしいでしょうか？')){
                $.ajax({
                    url: '/u/api/app/delete/' + util.getAppDetailId(),
                    type: 'delete'
                }).done((res) => {
                    location.href = '/'
                }).fail(() => {
                    alert("Error")
                })
            }
        })
    }
}