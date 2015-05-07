import config from './config'

var util = {

    getUserInfo() {
        var user = config.user
        if(!user || !user.ID) {
            return null
        }
        return user
    },

    getUrlParams() {
        var search = location.search.substring(1);
        return search ? JSON.parse('{"' + search.replace(/&/g, '","').replace(/=/g,'":"') + '"}',
                 function(key, value) { return key==="" ? value:decodeURIComponent(value) }) : {}
    },

    getImageHTML(id, w, h){
        // override by fb account id
        if(id.FBUser){
            id = id.FBUser.id
        }
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
            valid: function (value, query) {
                return true;
            },
            source:[{
                url: url,
                type:'remote',
                minLength: 1
            }],
            getTitle:function(item){
                return item['Name']
            },
            getValue:function(item){
                return item['Name']
            },
            render: function(item, source, pid, query){
                var id = item['ID'],
                    name = item['Name'];
                return '<div class="ma-friend-add-item' + (id == query ? ' active' : '') +
                '" data-id="' + encodeURIComponent(id) + '">' + util.getImageHTML(item) + '<span>' +  name + '</span></div>'
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
        if(url === "" || url === undefined){
            return true
        }
        var regex = new RegExp(/[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi);
        if(url.match(regex) ){
            return true
        }
        return false
    },

    isISODateFormat(dateStr) {
        if(dateStr === "" || dateStr === undefined){
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