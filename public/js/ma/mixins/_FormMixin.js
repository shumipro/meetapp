import $ from 'jquery'
import messages from '../messages'
import util from '../util'

export default class _FormMixin {
    constructor() {
    }

    validate(params) {
        // TODO: move to _FromMixin
        for(var prop in this.forms.props) {
            var info = this.forms.props[prop],
                value = params[prop]
            if(info.required){
                if($.trim(value) === ""){
                    return {"error": true, "message": "必須項目が入力されていません: " + (messages.props && messages.props[prop] || prop)}
                }
            }
            if(info.maxLength){
                if(value.length > info.maxLength){
                    var label = (messages.props && messages.props[prop] || prop) + "は" + info.maxLength + "文字以内で入力してください"
                    return {"error": true, "message": "入力された文字数が多すぎます: " + label}
                }
            }
            if(info.type === "url"){
                if(!util.isUrlFormat(value)){
                    return {"error": true, "message": "URLが不正です: " + (messages.props && messages.props[prop] || prop)}
                }
            }
            if(info.type === 'date'){
                if(!util.isISODateFormat(value)){
                    return {"error": true, "message": "日付が不正です: " + (messages.props && messages.props[prop] || prop)}
                }
            }
        }
        return {"error": false}
    }

    getParams() {
        // TODO: move to _FromMixin
        // collect params by forms DOM
        var param = {};
        for(var prop in this.forms.props) {
            var info = this.forms.props[prop]
            if(info.type === 'list'){
                var $inputs = $('*[data-list-name="' + prop + '"]')
                if($inputs.size() > 0) {
                    $inputs.each(function(index, target){
                        var $target = $(target),
                            obj = {}
                        if(!param[prop]){
                            param[prop] = []
                        }
                        if($target.prop('tagName').toLowerCase() === "input"){
                            var v = $target.val()
                            if(v !== ''){
                                obj[$target.attr('name')] = $target.val()
                                var listIndex = $target.data('list-index')
                                if(listIndex >= 0){
                                    param[prop][listIndex] = obj
                                }else{
                                    param[prop].push(obj)
                                }
                            }
                        }else{
                            // try to find children when the target itself does not have name attr
                            var $children = $target.find('*[name]')
                            $children.each(function(index, child){
                                var $child = $(child)
                                if($child.attr('type') == 'checkbox'){
                                    obj[$child.attr('name')] = $child.prop('checked')
                                }else{
                                    obj[$child.attr('name')] = $child.val()
                                }
                            });
                            param[prop].push(obj)
                        }
                    })
                }
            }else{
                var $input = $('input[name="' + prop + '"], select[name="' + prop + '"], textarea[name="' + prop + '"]')
                if($input.size() > 0) {
                    param[prop] = $input.val()
                }
            }
        }
        return param
    }
}