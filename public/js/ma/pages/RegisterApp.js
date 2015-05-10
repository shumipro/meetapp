import $ from 'jquery'
import autocomplete from 'jquery-autocomplete'
import util from '../util'
import config from '../config'
import constants from '../constants'
import _FormMixin from '../mixins/_FormMixin'
import ConstantSelect from '../components/ConstantSelect'
import Handlebars from 'handlebars'

var currentMemberEntryHtml = '<div class="ma-friend-add-member" data-list-name="currentMembers"><input type="hidden" name="id" value="{{id}}">' +
                      '{{{imgHtml}}}<span>{{name}}</span>' + 
                      '<select name="occupation" class="form-control ma-constant-select" data-constant="occupation"></select>' +
                      '<button type="button" class="btn btn-default ma-friend-delete-btn">削除</button>' +
                      '</div>'
var currentMemberEntryTmpl = Handlebars.compile(currentMemberEntryHtml)

var recruitMemberEntryHtml = '<div class="ma-friend-add-member" data-list-name="recruitMembers">' +
                      '<select name="occupation" class="form-control ma-constant-select" data-constant="occupation"></select>' +
                      '<button type="button" class="btn btn-default ma-friend-delete-btn">削除</button>' +
                      '</div>'
var recruitMemberEntryTmpl = Handlebars.compile(recruitMemberEntryHtml)

export default class RegisterApp extends _FormMixin {
    constructor() {
        super()
        this.forms = {
            props: {
                'name': { type: 'text', required: true, maxLength: 20 },
                'description': { type: 'text', maxLength: 1000 },
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
                'recruitMembers': {
                    type: 'list',
                    props: {
                        'occupation': { type: 'select' }
                    }
                },
                'demoUrl': { type: 'url' },
                'githubUrl': { type: 'url' },
                'meetingArea': { type: 'select' },
                'meetingFrequency': { type: 'select' },
                'projectStartDate': { type: 'date' },
                'projectReleaseDate': { type: 'date' }
            }
        }
        this._$submit = $('#ma_register_submitBtn')
        this._$submit.on('click', this.submit.bind(this))

        // auto complete setup
        util.autoCompleteAddInit("/api/user/search/%QUERY%", $('#ma_register_add_currentMember_suggest_input'), $('#ma_register_add_currentMember_suggest_btn'), (item) => {
            // check currentMembers already has the member
            var members = this.getParams().currentMembers;
            for(var i=0; i<members.length; i++){
                if(members[i].id === item.ID){
                    alert("そのユーザーは既に追加されています")
                    return;
                }
            }
            this.createCurrentMemberEntry(item)
        })
        $('#ma_register_add_recruitMember_btn').on('click', () => {
            this.createRecruitMemberEntry()
        })
        // attach delete button
        var $deleteBtn = $('.ma-friend-delete-btn')
        $deleteBtn.on('click', function() {
            $(this).parent().remove()
        })
    }

    createCurrentMemberEntry(item) {
        var $item = $(currentMemberEntryTmpl({
            id: item.ID,
            name: item.Name,
            imgHtml: util.getImageHTML(item)
        }))
        this._createMemberEntry($item, $('#ma_register_add_currentMember_result'))
    }

    createRecruitMemberEntry() {
        var $item = $(recruitMemberEntryTmpl({}))
        this._createMemberEntry($item, $('#ma_register_add_recruitMember_result'))
    }

    _createMemberEntry($item, $wrap) {
        // attach events
        var $deleteBtn = $item.find('.ma-friend-delete-btn')
        $deleteBtn.on('click', () => {
            $deleteBtn.parent().remove()
        })
        // create occupations select
        new ConstantSelect($item.find('.ma-constant-select'))
        $wrap.append($item)
    }

    submit() {
        var params = this.getParams()
        //  validation
        var result = this.validate(params)
        if(result.error){
            alert(result.message)
            return
        }
        var appId = this._$submit.data("app-id")
        if(appId){
            params.id = appId
        }
        $.ajax({
            url: appId ? '/u/api/app/edit/' + appId : '/u/api/app/register',
            type: appId ? 'put' : 'post',
            contentType:"application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify(params)
        }).done((res) => {
            if(res && res.id){
                location.href = '/app/detail/' + res.id
            }
        }).fail(() => {
            alert("Error")
        })
    }
}
