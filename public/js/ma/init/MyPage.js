import $ from 'jquery'
import _InitBase from './_InitBase'
import ImageUploader from '../components/ImageUploader'

export default class MyPage extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new ImageUploader('/u/api/upload/image')
    }
}