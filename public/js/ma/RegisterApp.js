import $ from 'jquery'

var FORM_PREFIX = 'ma_register_form_'

export default class RegisterApp {
    constructor() {
        var form = this.forms = {
            props: {
                'name': { type: 'text' },
                'description': { type: 'text' },
                'platform': { type: 'select' },
                'category': { type: 'select' },
                'pLang': { type: 'select' },
                'keywords': { type: 'text' },
                'images': {
                    type: 'list',
                    props: {
                        'url': {type: 'url' }
                    }
                },
                'currentMembers': {
                    type: 'list',
                    props: {
                        'name': { type: 'text' },
                        'role': { type: 'select' }
                    }
                },
                'wantMembers': {
                    type: 'list',
                    props: {
                        'role': { type: 'select' }
                    }
                },
                'demoUrl': { type: 'url' },
                'githubUrl': { type: 'url' },
                'meetingPlace': { type: 'text' },
                'meetingOften': { type: 'text' },
                'projectStartDate': { type: 'date' },
                'projectReleaseDate': { type: 'date' }
            }
        }
        this.$addCurrentMember = $('#ma_register_add_currentMember')
        this.$addWantMember = $('#ma_register_add_wantMember')
        this.$submitBtn = $('#ma_register_submitBtn')
        this.$submitBtn.on('click', this.submit.bind(this))
    }

    submit() {
        // {"name": "App name", "description": "hoge", "images": [{"url": "https://golang.org/doc/gopher/gopherbw.png"}]}
       $.ajax({
            url: '/api/app/register',
            type: 'post',
            contentType:"application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify(this.getParams())
        }).done((res) => {
            if(res && res.id){
                location.href = '/app/detail/' + res.id
            }
        }).fail(() => {
            alert("Error")
        })
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
                            v = $target.val(),
                            obj = {}
                        if(v !== ''){
                            if(!param[prop]){
                                param[prop] = []
                            }
                            obj[$target.attr('name')] = $target.val()
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