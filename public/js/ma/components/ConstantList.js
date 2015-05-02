import $ from 'jquery'
import constants from '../constants'
import Handlebars from 'handlebars'

var html = '<div class="col-md-2 col-sm-2 ma-app-tag">' +
            '<a href="/app/list?{{prop}}={{id}}">' + 
            '<h4 class="media-heading">{{name}}</h4>' +
            '</a>'
            '</div>'
var tmpl = Handlebars.compile(html)

var footerHtml = '<li><a href="/app/list?{{prop}}={{id}}">{{name}}</a></li>'
var footerTmpl = Handlebars.compile(footerHtml)

export default class ConstantList {
    constructor() {
        // for top page
        // e.g. <div class="ma-constant-list" data-constant="platforms"></div>
        var $lists = $('.ma-constant-list')
        $lists.each(function(index, target){
            var $widget = $(target),
                prop = $widget.data('constant')
            if(prop && constants[prop]){
                var list = constants[prop]
                for(var i = 0; i < list.length; i++) {
                    $(tmpl({
                        id: list[i].id,
                        name: list[i].name,
                        prop: prop
                    })).appendTo($widget)
                }
            }
        })
        // for footer
        // e.g. <ul class="ma-constant-footer-list" data-constant="occupations"></ul>
        var $lists = $('.ma-constant-footer-list')
        $lists.each(function(index, target){
            var $widget = $(target),
                prop = $widget.data('constant')
            if(prop && constants[prop]){
                var list = constants[prop]
                for(var i = 0; i < 8; i++) {
                    $(footerTmpl({
                        id: list[i].id,
                        name: list[i].name,
                        prop: prop
                    })).appendTo($widget)
                }
                $widget.append($('<li><a href="/app/list">もっと見る</a></li>'))
            }
        })
    }
}