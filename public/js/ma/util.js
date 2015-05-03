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
    },

    isUrlFormat(url) {
        if(url === ""){
            return true
        }
        var regex = new RegExp(/[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi);
        if(url.match(regex) ){
            return true
        }
        return false
    },

    isISODateFormat(dateStr) {
        if(dateStr === ""){
            return true
        }
        if(dateStr.length === 10 && dateStr.match(/(\d{4})-(\d{2})-(\d{2})/)){
            return true
        }
        return false
    },

    getAppDetailId(){
        // get id from current URL
        var result = location.pathname.match(new RegExp(/\/app\/detail\/([\w-]+)/))
        return result && result[1] || ""
    }
}

module.exports = util