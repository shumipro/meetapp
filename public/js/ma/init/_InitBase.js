import ConstantList from '../components/ConstantList'
import AppButtons from '../components/AppButtons'

export default class _InitBase {
    constructor() {
        
    }

    init() {
        new ConstantList()
        new AppButtons()
    }
}