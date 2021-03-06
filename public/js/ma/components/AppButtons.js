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

        // 興味ありボタン
        $('.ma-join-btn').on('click', () => {
            var _requesting = false
            return () => {
                if(util.getUserInfo()){
                    if(_requesting){
                        return
                    }
                    $.ajax({
                        url: '/u/api/app/join/' + util.getAppDetailId(),
                        type: 'post'
                    }).done((res) => {
                        alert('開発アイデアに興味をもっていただきありがとうございます。管理者の方からのご連絡をお待ち下さい。')
                        _requesting = false
                    }).fail(() => {
                        alert("Error")
                        _requesting = false
                    })            
                    _requesting = true
                }else{
                    // move to login for anonymous
                    location.href = "/login"
                }
            }
        }())
    }
}