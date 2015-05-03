import $ from 'jquery'

var SELECTOR = '.ma-github-issues'

export default class GihubIssues {
    constructor() {
        var $widget = $(SELECTOR)
        if($widget.size() > 0){
            var url = $widget.data('url')
            if(url){
                // get from github URL e.g. https://github.com/shumipro/meetapp
                var result = url.match(new RegExp(/https:\/\/github.com\/([\w-]+)\/([\w-]+)/))
                if(result && result.length > 2) {
                    var api = 'https://api.github.com/repos/' + result[1] + '/' + result[2] + '/issues?'
                    $.ajax({
                        url: api,
                        dataType: 'jsonp'
                    }).done((res) => {
                        this.render(res)
                    })
                }
            }
        }
    }

    render(results) {
        var $widget = $(SELECTOR)
        var resultDataList = results.data
        var limit = 10
        if(!resultDataList || resultDataList.length === 0){
            $widget.append('<div>現在Github Issueはありません</div>')
            return
        }
        var $ul = $('<ul></ul>')
        for (var i = 0; i < limit; i++) {
            var issue = resultDataList[i]
            if(issue) {
                var $li = $('<li></li>')
                var $a = $('<a class="ma-github-issue-link" href="' + issue.html_url + '" target="_blank"><span class="ma-github-issue-icon"></span>' + issue.title + '</a>')
                $li.append($a)
                $ul.append($li)
            }
        }
        $widget.append($ul)
    }
}