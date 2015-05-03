import _InitBase from './_InitBase'
import Discussion from '../components/Discussion'
import StarButton from '../components/StarButton'
import GithubIssues from '../components/GithubIssues'

export default class Detail extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new GithubIssues()
        new Discussion()
    }
}