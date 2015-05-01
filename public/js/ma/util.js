var util = {

    loadJSONP(api, callbackName) {
        var head = document.getElementsByTagName('head')[0];
        var el = document.createElement('script');
        el.src = api + '&callback=' + callbackName;
        head.insertBefore(el, head.firstChild);
    },

    getImageHTML(id, w, h){
        w = w || "32px", h = h || "32px";
        return "<img class='img-rounded' width='" + w + "' height='" + h + "' src='" + util.getImageURL(id) + "'/>";
    },

    getImageURL(id){
        return "http://graph.facebook.com/" + id + "/picture?type=square";
    },

    autoCompleteAddInit(url, $input, $addBtn, addCallback){
        var _selectedItem = "";
        $input.autocomplete({
            openOnFocus: false,
            appendMethod:'replace',
            valid: function () {
                return true;
            },
            source:[{
                url: url,
                type:'remote'
            }],
            getTitle:function(item){
                return item['name']
            },
            getValue:function(item){
                return item['name']
            },
            render: function(item, source, pid, query){
                var id = item['id'],
                    name = item['name'];
                return '<div class="ma-friend-add-item' + (id == query ? ' active' : '') +
                '" data-id="' + encodeURIComponent(id) + '">' + util.getImageHTML(id) + '<span>' +  name + '</span></div>'
            }
        }).on('selected.xdsoft',function(e, data){
            _selectedItem = data;
        });
        $addBtn.click(function(){
            if(_selectedItem){
                if(addCallback){ addCallback(_selectedItem) }
                $input.val("");
                _selectedItem = null;
            }
        });
    }   
}

module.exports = util