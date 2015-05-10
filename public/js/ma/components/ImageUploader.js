import $ from 'jquery'

export default class ImageUploader {
    constructor(api, selector) {
        this._api = api
        this._$file = $(selector)
    }

    upload(){
        var formData = new FormData();
        $.each(this._$file[0].files, function(i, file){
            formData.append('file', file)
        })
        return $.ajax({
            url: this._api,
            method: 'post',
            dataType: 'json',
            data: formData,
            processData: false,
            contentType: false
        })
    }

    validate(){
        // TODO: size valdation
    }

    getFileName(){
      return this._$file.val().replace(/\\/g, '/').replace(/.*\//, '')
    }
}