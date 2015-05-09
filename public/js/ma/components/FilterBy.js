import $ from 'jquery'
import messages from '../messages'
import constants from '../constants'
import util from '../util'
import Handlebars from 'handlebars'

var html = '<div class="ma-tag">{{propLabel}}: {{name}}</div>'
var tmpl = Handlebars.compile(html)

export default class ConstantFilter {
    constructor() {
        var $wrap = $('.ma-filter-by-wrap'),
            params = util.getUrlParams(),
            count = 0
        for(var prop in params){
            if(prop && params[prop] && constants[prop]){
                // get label from id
                var $item = $(tmpl({
                    name: util.getConstantLabel(prop, params[prop]),
                    propLabel: messages.constants && messages.constants[prop] || prop,
                }))
                $wrap.append($item)
                count++
            }
        }
        if(count > 0){
            $wrap.show()
        }
    }
}