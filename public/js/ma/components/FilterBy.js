import $ from 'jquery'
import constants from '../constants'
import util from '../util'
import Handlebars from 'handlebars'

var html = '<span>Filter {{prop}}: {{name}}</span>'
var tmpl = Handlebars.compile(html)

export default class ConstantFilter {
    constructor() {
    //     var $wrap = $('.ma-filter-by-wrap'),
    //         params = util.getUrlParams()
    //     for(var prop in params){
    //         if(prop && params[prop]){
    //             var id = params[prop]
    //             // TODO: get label from id
    //             var $item = $(tmpl({
    //                 name: id,
    //                 prop: prop
    //             }))
    //             $wrap.append($item)
    //         }
    //     }
    }
}