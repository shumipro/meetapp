import _InitBase from './_InitBase'
import ConstantFilter from '../components/ConstantFilter'

export default class Register extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new ConstantFilter()
    }
}