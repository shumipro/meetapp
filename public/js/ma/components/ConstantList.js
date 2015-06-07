import $ from 'jquery'
import constants from '../constants'
import config from '../config'
import Handlebars from 'handlebars'

var listHtml = '<div class="col-md-2 col-sm-2 ma-app-tag">' +
            '<a href="/app/list?{{prop}}={{id}}">' + 
            '<h4 class="media-heading">{{name}}</h4>' +
            '</a>'
            '</div>'
var listTmpl = Handlebars.compile(listHtml)

var iconListHtml = '<div class="ma-constant-icon-list-item">' +
            '<a href="/app/list?{{prop}}={{id}}">' + 
            '<div><img src="' + config.static_path + 'img/{{prop}}/{{id}}.png"></div>' + 
            '<div class="ma-constant-icon-list-title">{{name}}</div>' +
            '</a>'
            '</div>'
var iconListTmpl = Handlebars.compile(iconListHtml)


var footerHtml = '<li><a href="/app/list?{{prop}}={{id}}">{{name}}</a></li>'
var footerTmpl = Handlebars.compile(footerHtml)
var FOOTER_ITEM_COUNT = 8

export default class ConstantList {
    constructor() {
        // for top page
        // e.g. <div class="ma-constant-list" data-constant="platform"></div>
        var $lists = $('.ma-constant-list')
        $lists.each((index, target) => {
            this.render($(target), listTmpl)
        })
        // list with icons
        // e.g. <div class="ma-constant-icon-list" data-constant="occupation"></div>
        var $lists = $('.ma-constant-icon-list')
        $lists.each((index, target) => {
            this.render($(target), iconListTmpl)
        })
        // for footer
        // e.g. <ul class="ma-constant-footer-list" data-constant="occupation"></ul>
        var $lists = $('.ma-constant-footer-list')
        $lists.each((index, target) => {
            var $widget = $(target)
            this.render($widget, footerTmpl, FOOTER_ITEM_COUNT)
            $widget.append($('<li><a href="/app/list">もっと見る</a></li>'))
        })
    }

    render($widget, tmpl, maxCount) {
        var prop = $widget.data('constant')
        if(prop && constants[prop]){
            var list = this._filterListByValues(constants[prop], $widget.data('values')),
                len = maxCount || list.length
            for(var i = 0; i < len; i++) {
                $(tmpl({
                    id: list[i].id,
                    name: list[i].name,
                    prop: prop
                })).appendTo($widget)
            }
        }
    }

    /**
     * Return items by comma separated values
     * e.g. data-values="0,4,3,1"
     * list: [{"id":"0","name":"まだ決めていない"},{"id":"1","name":"渋谷"}...]
     */
    _filterListByValues(list, valuesStr) {
        if(valuesStr === undefined || valuesStr === ""){
            return list
        }
        var arr = [],
            values = valuesStr.split(',')
        values.forEach((value) => {
            // for last comma
            if(value !== ''){
                for(var i=0; i<list.length; i++){
                    if(list[i].id === value){
                        arr.push(list[i])
                        break
                    }
                }
            }
        })
        return arr
    }
}