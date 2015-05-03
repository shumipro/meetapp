import $ from 'jquery'
import constants from '../constants'
import Handlebars from 'handlebars'

var html = '<li><a href="/app/list?{{prop}}={{id}}">{{name}}</a></li>'
var tmpl = Handlebars.compile(html)

export default class ConstantFilter {
    constructor() {
        // for top page
        // e.g. <div class="ma-constant-list" data-constant="platforms"></div>
        var $lists = $('.ma-constant-filter')
        $lists.each(function(index, target){
            var $widget = $(target),
                prop = $widget.data('constant')
            if(prop && constants[prop]){
                var list = constants[prop]
                for(var i = 0; i < list.length; i++) {
                    var $item = $(tmpl({
                        id: list[i].id,
                        name: list[i].name,
                        prop: prop
                    }))
                    // TODO: add param and move to page
                    // $item.on('click', () => {})
                    $item.appendTo($widget)
                }
            }
        })
    }
}