import $ from 'jquery'
import constants from '../constants'

export default class ConstantSelect {
    constructor() {
        // setup options by using constants
        // e.g. <select name="platform" class="ma-constant-select" data-constant="platforms"></select>
        var $selects = $('.ma-constant-select')
        $selects.each(function(index, select){
            var $widget = $(select),
                prop = $widget.data('constant')
            if(prop && constants[prop]){
                var options = [],
                    list = constants[prop]
                for(var i = 0; i < list.length; i++) {
                    options.push("<option value='" + list[i].id + "'>" + list[i].name + "</option>")
                }
                $widget.append(options.join(''))
            }
        })
    }
}