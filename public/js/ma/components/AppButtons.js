import $ from 'jquery'
import util from '../util'

export default class ConstantSelect {
    constructor() {
        // common button widgets
        $('.ma-app-edit-btn').on('click', () => {
            // TODO: move to edit page

        })

        $('.ma-app-delete-btn').on('click', () => {
            if(window.confirm('この開発アイデアを削除してもよろしいでしょうか？')){
                // TODO: delete
                alert(util.getAppDetailId())
            }
        })
        // TODO: star button
        $('.ma-app-star-btn').on('click', () => {
            // check the user is already logged in

            // move to login screen

            // send star
        })

        // set current URL for share buttons
        $('.fb-like[data-href=""]').data('href', location.href)
        $('.twitter-share-button[data-url=""]').data('href', location.href)
    }
}