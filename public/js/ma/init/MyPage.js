import $ from 'jquery'
import _InitBase from './_InitBase'
import MypageEdit from '../pages/MypageEdit'

export default class MyPage extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new MypageEdit()
    }
}