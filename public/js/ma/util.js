module.exports = {

    loadJSONP(api, callbackName) {
        var head = document.getElementsByTagName('head')[0];
        var el = document.createElement('script');
        el.src = api + '&callback=' + callbackName;
        head.insertBefore(el, head.firstChild);
    }

}