import $ from 'jquery'

export default class ConstantSelect {
    constructor() {
        // common button widgets
        $('.ma-app-edit-btn').on('click', () => {
            // TODO: move to edit page

        })

        $('.ma-app-delete-btn').on('click', () => {
            if(window.confirm('この開発アイデアを削除してもよろしいでしょうか？')){
                // TODO: delete
            }
        })
    }
}