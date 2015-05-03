import $ from 'jquery'
import constants from '../constants'

export default class ConstantSelect {
    constructor($selects) {
        // setup options by using constants
        // e.g. <select name="platform" class="ma-constant-select" data-constant="platforms"></select>
        $selects.each(function(index, select){
            var $widget = $(select),
                prop = $widget.data('constant')
            if(prop && constants[prop]){
                var options = [],
                    list = constants[prop]
                var selectedValue = null
                for(var i = 0; i < list.length; i++) {
                    var value = list[i].id,
                        defaultValue = $widget.data('default-value')
                    if(value === defaultValue + ""){
                        selectedValue = defaultValue
                    }
                    options.push("<option value='" + value + "' >" + list[i].name + "</option>")
                }
                $widget.append(options.join(''))
                if(selectedValue){
                    $widget.val(selectedValue)
                }
            }
        })
    }
}