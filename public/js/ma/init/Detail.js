import _InitBase from './_InitBase'
import StarButton from '../components/StarButton'
import Discussion from '../components/Discussion'
import GithubIssues from '../components/GithubIssues'

export default class Detail extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new StarButton()
        new GithubIssues()
        new Discussion()
    }
}