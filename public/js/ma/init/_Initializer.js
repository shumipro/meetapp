import _InitBase from './_InitBase'
import Detail from './Detail'
import List from './List'
import Register from './Register'

module.exports = {

    // define a map URL path to initializer
    routeMap: {
        '/html/detail.html': Detail,
        '/app/detail/:id': Detail,
        '/html/list.html': List,
        '/app/list': List,
        '/html/register.html': Register,
        '/app/register': Register
    },

    getInitializer() {
        // set default
        var Clazz = _InitBase,
            path = location.pathname,
            keys = Object.keys(this.routeMap)
        // iterate to match route definitions
        for(var i = 0; i < keys.length; i++){
            var route = keys[i]
            if(this.isMatch(route, path)) {
                Clazz = this.routeMap[route]
                break
            }
        }
        // instanciate initializer
        return new Clazz()
    },

    isMatch(route, path) {
        // reaplace :xxx
        var routeMatcher = new RegExp(route.replace(/:[^\s/]+/g, '([\\w-]+)'))
        return path.match(routeMatcher)
    }

}