import $ from 'jquery'
import autocomplete from 'jquery-autocomplete'
import util from '../util'
import config from '../config'
import constants from '../constants'
import Handlebars from 'handlebars'

var FORM_PREFIX = 'ma_register_form_'

var currentMemberEntryHtml = '<div class="ma-friend-add-member" data-list-name="currentMembers"><input type="hidden" name="id" value="{{id}}">' +
                      '{{{imgHtml}}}<span>{{name}}</span>' + 
                      '<select class="form-control" name="occupation">{{#each occupations}}<option value="{{id}}">{{name}}</option>{{/each}}</select>' +
                      '<button type="button" class="btn btn-default">削除</button>'
                      '</div>'
var currentMemberEntryTmpl = Handlebars.compile(currentMemberEntryHtml)

var wantMemberEntryHtml = '<div class="ma-friend-add-member" data-list-name="wantMembers">' +
                      '<select class="form-control" name="occupation">{{#each occupations}}<option value="{{id}}">{{name}}</option>{{/each}}</select>' +
                      '<button type="button" class="btn btn-default">削除</button>'
                      '</div>'
var wantMemberEntryTmpl = Handlebars.compile(wantMemberEntryHtml)

export default class RegisterApp {
    constructor() {
        var form = this.forms = {
            props: {
                'name': { type: 'text', required: true, maxLength: 20 },
                'description': { type: 'text', maxLength: 400 },
                'platform': { type: 'select' },
                'category': { type: 'select' },
                'pLang': { type: 'select' },
                'keywords': { type: 'text', maxLength: 100 },
                'images': {
                    type: 'list',
                    props: {
                        'url': {type: 'url' }
                    }
                },
                'currentMembers': {
                    type: 'list',
                    props: {
                        'id': { type: 'text' },
                        'occupation': { type: 'select' }
                    }
                },
                'wantMembers': {
                    type: 'list',
                    props: {
                        'occupation': { type: 'select' }
                    }
                },
                'demoUrl': { type: 'url' },
                'githubUrl': { type: 'url' },
                'meetingArea': { type: 'select' },
                'meetingOften': { type: 'select' },
                'projectStartDate': { type: 'date' },
                'projectReleaseDate': { type: 'date' }
            }
        }
        this.$submitBtn = $('#ma_register_submitBtn')
        this.$submitBtn.on('click', this.submit.bind(this))

        // auto complete setup
        util.autoCompleteAddInit("/api/search_friend?q=%QUERY%", $('#ma_register_add_currentMember_suggest_input'), $('#ma_register_add_currentMember_suggest_btn'), (item) => {
            // TODO: check currentMembers already has the member
            this.createCurrentMemberEntry(item)
        })
        this.$addWantMember = $('#ma_register_add_wantMember_btn')
        this.$addWantMember.on("click", () => {
            this.createWantMemberEntry()
        })
        // set myself at init
        this.createCurrentMemberEntry(config.user)
        // TODO: set values for edit
        // var override = {
        //     "name": "App name",
        //     "description": "hoge",
        //     "platform": "3",
        //     "category": "4",
        //     "pLang": "5",
        //     "keywords": "keyword test",
        //     "images": [{"url": "https://golang.org/doc/gopher/gopherbw.png"}],
        //     "currentMembers": [
        //         {id: "1234", name: "Tejitak", occupation: "4"},
        //         {id: "1234", name: "kyokomiさん", occupation: "2"}
        //     ],
        //     "wantMembers": [
        //         {occupation: "5"},
        //         {occupation: "3"},
        //     ],
        //     "demoUrl": "http://demo.com/",
        //     "githubUrl": "http://github.com/",
        //     "meetingArea": "1",
        //     "meetingOften": "1",
        //     "projectStartDate": "2015-04-30",
        //     "projectReleaseDate": "2015-05-30"
        // }
        // this.update(override)
    }

    createCurrentMemberEntry(item) {
        var $item = $(currentMemberEntryTmpl({
            id: item.id,
            name: item.name,
            imgHtml: util.getImageHTML(item.id),
            occupations: constants.occupation
        }))
        this._createMemberEntry($item, $('#ma_register_add_currentMember_result'))
    }

    createWantMemberEntry() {
        var $item = $(wantMemberEntryTmpl({
            occupations: constants.occupation
        }))
        this._createMemberEntry($item, $('#ma_register_add_wantMember_result'))
    }

    _createMemberEntry($item, $wrap) {
        var $deleteBtn = $item.find('button')
        $deleteBtn.on("click", () => {
            $deleteBtn.parent().remove()
        })
        $wrap.append($item)
    }

    update(params) {
        // update DOM by given parmas
        
    }

    submit() {
        var params = this.getParams()
        //  validation
        var result = this.validate(params)
        if(result.error){
            alert(result.message)
            return
        }
        // {"name": "App name", "description": "hoge", "images": [{"url": "https://golang.org/doc/gopher/gopherbw.png"}]}
        $.ajax({
            url: '/api/app/register',
            type: 'post',
            contentType:"application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify(params)
        }).done((res) => {
            if(res && res.ID){
                location.href = '/app/detail/' + res.ID
            }
        }).fail(() => {
            alert("Error")
        })
    }

    validate(params) {
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

        }
        return {"error": false}
    }

    getParams() {
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
                                param[prop].push(obj)
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
                var $input = $('*[name="' + prop + '"]')
                if($input.size() > 0) {
                    param[prop] = $input.val()
                }
            }
        }
        return param
    }

}