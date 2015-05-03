import _InitBase from './_InitBase'
import GithubIssues from '../components/GithubIssues'

export default class Detail extends _InitBase {
    constructor() {
        super()
    }

    init() {
        super.init()
        new GithubIssues()
    }
}