import _InitBase from './_InitBase'
import ConstantFilter from '../components/ConstantFilter'
import FilterBy from '../components/FilterBy'

export default class Register extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new ConstantFilter()
        new FilterBy()
    }
}