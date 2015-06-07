import ConstantList from '../components/ConstantList'
import AppButtons from '../components/AppButtons'
import Notifications from '../components/Notifications'

export default class _InitBase {
    constructor() {
        
    }

    init() {
        new ConstantList()
        new AppButtons()
        new Notifications()
    }
}