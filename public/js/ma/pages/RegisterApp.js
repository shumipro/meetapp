import $ from 'jquery'
import autocomplete from 'jquery-autocomplete'
import util from '../util'
import config from '../config'
import constants from '../constants'
import _FormMixin from '../mixins/_FormMixin'
import ConstantSelect from '../components/ConstantSelect'
import ImageUploader from '../components/ImageUploader'
import Handlebars from 'handlebars'

var currentMemberEntryHtml = '<div class="ma-friend-add-member" data-list-name="currentMembers"><input type="hidden" name="id" value="{{id}}">' +
                      '{{{imgHtml}}}<span>{{name}}</span>' + 
                      '<select name="occupation" class="form-control ma-constant-select" data-constant="occupation"></select>' +
                      '<label><input type="checkbox" class="form-control" name="isAdmin">管理者</label>' +
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
                'name': { type: 'text', required: true, maxLength: 30 },
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
                        'occupation': { type: 'select' },
                        'isAdmin': { type: 'checkbox'}
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
        // setup image uploaders
        this._uploaders = []
        $('.ma-register-form-image-file').each((index, elm) => {
            this._uploaders.push(new ImageUploader('/u/api/upload/image', elm))
        })
        // delete registered images
        $('.register-form-image-delete-btn').on('click', function(){
            var $btn = $(this)
            $btn.parent().parent().find('input[type="text"]').val('')
            $btn.hide()
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
        var result = this.validateImages()
        if(result.error){
            return
        }

        var uploadList = [],
            deferredList = []

        // remove empty object from list
        $('.ma-register-form-image-input').each((index, elm) => {
            var $input = $(elm),
                v = $input.val()
            if(v){
                var uploader = this._uploaders[index]
                if(v.indexOf('http://') === 0 || v.indexOf('https://') === 0){
                    // keep url already set
                    uploader.dummy = true
                }
                uploadList.push(uploader)
            }
        })
        uploadList.forEach((uploader, i) => {
            var $dfd = $.Deferred()
            if(uploader.dummy){
                $dfd.resolve()
            }else{
                uploader.upload("apps").then((res) => {
                    // override images
                    params.images[i].url = res.ImageURL
                    $dfd.resolve(res)
                })
            }
            deferredList.push($dfd)
        })
        if(this._requesting){
            return
        }
        this._requesting = true
        this._orgSubmitLabel = this._$submit.val()
        this._$submit.val('保存中...')
        if(deferredList.length > 0){
            // convert array to func arguments
            $.when.apply($, deferredList).then(() => {
                this.postApp(params)
            })
        }else{
            this.postApp(params)
        }
    }

    validateImages(){
        // image file size validate
        for(var i=0; i<this._uploaders.length; i++){
            if(!this._uploaders[i].validate()){
                return {error: true}
            }
        }
        return {error: false}
    }

    postApp(params){
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
            this._requesting = false
        }).fail(() => {
            alert("Error")
            this._requesting = false
            this._$submit.val(this._orgSubmitLabel)
        })
    }
}
