import _InitBase from './_InitBase'
import ConstantFilter from '../components/ConstantFilter'
import FilterBy from '../components/FilterBy'
import ListPaging from '../components/ListPaging'

export default class Register extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new ConstantFilter()
        new FilterBy()
        new ListPaging()
    }
}