import $ from 'jquery'
import util from '../util'
import _FormMixin from '../mixins/_FormMixin'
import ImageUploader from '../components/ImageUploader'

export default class MypageEdit extends _FormMixin {
    constructor() {
        super()
        this._uploader = new ImageUploader('/u/api/upload/image', '#ma-mypage-edit-profile-image-file')
        this.forms = {
            props: {
                'Name': { type: 'text', required: true, maxLength: 20 },
                'Comment': { type: 'text', maxLength: 50 },
                'HomePageURL': { type: 'url' },
                'GitHubURL': { type: 'url' },
            }
        }
        this._$submit = $('.ma-profile-edit-update-btn')
        this._$submit.on('click', this.submit.bind(this))
    }

    submit() {
        var params = this.getParams()
        //  validation
        var result = this.validate(params)
        if(result.error){
            alert(result.message)
            return
        }
        // set user ID
        var user = util.getUserInfo()
        if(!user){ return }
        params.ID = user.ID
        params.ImageURL = user.ImageURL
        params.LargeImageURL = user.LargeImageURL
        if(this._uploader.getFileName() && this._uploader.validate()){
            // upload image
            this._uploader.upload("users").then((res) => {
                // override images
                params.ImageURL = res.ImageURL
                params.LargeImageURL = res.LargeImageURL
                this.postMypage(params)            
            }, () => {
                alert('Upload error')
            })
        }else{
            // imageがセットされていない時はスキップ
            this.postMypage(params)
        }
    }

    postMypage(params) {
        $.ajax({
            url: '/u/api/user',
            type: 'put',
            contentType:"application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify(params)
        }).done((res) => {
            if(res){
                location.href = '/u/mypage'
            }
        }).fail(() => {
            alert("Error")
        })
    }
}