import $ from 'jquery'
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
                    return {"error": true, "message": "必須項目が入力されていません: " + prop}
                }
            }
            if(info.maxLength){
                if(value.length > info.maxLength){
                    return {"error": true, "message": "入力されたデータが長すぎます: " + prop}
                }
            }
            if(info.type === "url"){
                if(!util.isUrlFormat(value)){
                    return {"error": true, "message": "URLが不正です: " + prop}
                }
            }
            if(info.type === 'date'){
                if(!util.isISODateFormat(value)){
                    return {"error": true, "message": "日付が不正です: " + prop}
                }
            }
            // temp
            if(prop === "images"){
                for(var i=0; i<value.length; i++){
                    var img = value[i]
                    if(!util.isUrlFormat(img.url)){
                        return {"error": true, "message": "URLが不正です: " + prop}
                    }
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
                                obj[$child.attr('name')] = $child.val()
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