import _InitBase from './_InitBase'
import RegisterApp from '../components/RegisterApp'
import ConstantSelect from '../components/ConstantSelect'

export default class Register extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new RegisterApp()
        new ConstantSelect()
    }
}