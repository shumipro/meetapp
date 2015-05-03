import $ from 'jquery'
import util from '../util'

export default class AppButtons {
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
        // set current URL for share buttons
        $('.fb-like[data-href=""]').data('href', location.href)
        $('.twitter-share-button[data-url=""]').data('href', location.href)
    }
}