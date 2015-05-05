import $ from 'jquery'
import constants from '../constants'
import util from '../util'
import Handlebars from 'handlebars'

var html = '<li><a href="javascript:;">{{name}}</a></li>'
var tmpl = Handlebars.compile(html)

export default class ConstantFilter {
    constructor() {
        // for top page
        // e.g. <div class="ma-constant-list" data-constant="platforms"></div>
        var $lists = $('.ma-constant-filter'),
            params = util.getUrlParams()
        $lists.each(function(index, target){
            var $widget = $(target),
                prop = $widget.data('constant')
            if(prop && constants[prop]){
                var list = constants[prop]
                $.each(list, function(i, item){
                    var value = item.id
                    var $item = $(tmpl({
                        id: value,
                        name: item.name,
                        prop: prop
                    }))
                    // set title
                    $item.attr('title', item.name + 'で絞り込む')
                    // 既に絞り込まれている値に関してはハイライト表示
                    if(params[prop] && params[prop] === value){
                        $item.addClass('ExSelected')
                        // title上書き
                        $item.attr('title', item.name + 'の絞り込みを解除')
                    }
                    // add param and move to page /app/list?{{prop}}={{id}}
                    $item.on('click', () => {
                        // filterで絞り込んだときはpageのパラメーターを削除する
                        delete params['page']
                        // 現在の絞り込みを解除する
                        delete params[prop]
                        if(!$item.hasClass('ExSelected')){
                            params[prop] = value
                        }
                        location.href = '/app/list?' + $.param(params)
                    })
                    $item.appendTo($widget)
                });
            }
        })
    }
}