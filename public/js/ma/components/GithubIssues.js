import $ from 'jquery'
import util from '../util'

var SELECTOR = '.ma-github-issues'

// handle jsonp in global scope
window.APICallback = (results) => {
    var $widget = $(SELECTOR)
    var resultDataList = results.data
    var $ul = $('<ul></ul>')
    var limit = resultDataList.length
    for (var i = 0; i < limit; i++) {
        var issue = resultDataList[i]
        var $li = $('<li></li>')
        var $a = $('<a class="ma-github-issue-link" href="' + issue.html_url + '" target="_blank"><span class="ma-github-issue-icon"></span>' + issue.title + '</a>')
        $li.append($a)
        $ul.append($li)
    }
    $widget.append($ul)
}

export default class GihubIssues {
    constructor() {
        var $widget = $(SELECTOR)
        if($widget.size() > 0){
            var url = $widget.data('url')
            if(url){
                // TODO: get from github URL
                var owner = 'shumipro'
                var repo = 'meetapp'
                var api = 'https://api.github.com/repos/' + owner + '/' + repo + '/issues?'
                util.loadJSONP(api, 'APICallback');
            }
        }
    }

}