import $ from 'jquery'

export default class ConstantFilter {
    constructor(api) {
        var $uploadBtn = $('.ma-image-upload-btn')
        $uploadBtn.on('click', function() {
            var form = $('.ma-image-upload-form').get()[0];
            var formData = new FormData( form );
            $.ajax({
                url: api,
                method: 'post',
                dataType: 'json',
                data: formData,
                processData: false,
                contentType: false
            }).done(function( res ) {
                console.log( 'SUCCESS', res )
                location.reload()
            }).fail(function( jqXHR, textStatus, errorThrown ) {
                console.log( 'ERROR', jqXHR, textStatus, errorThrown )
            })
            return false
        })
    }
}