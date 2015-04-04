import $ from 'jquery'
import Bootstrap from 'bootstrap/dist/js/npm'
import BootstrapExtend from './BootstrapExtend'
import ScrollEffect from './ScrollEffect'

$(document).ready(function() {
    (new ScrollEffect()).init()
})