import $ from 'jquery'
import _InitBase from './_InitBase'
import RegisterApp from '../pages/RegisterApp'
import ConstantSelect from '../components/ConstantSelect'

export default class Register extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new RegisterApp()
        new ConstantSelect($('.ma-constant-select'))
    }
}