import $ from 'jquery'
import Bootstrap from 'bootstrap/dist/js/npm'
import ScrollEffect from './ScrollEffect'
import RegisterApp from './RegisterApp'

$(document).ready(() => {
    new ScrollEffect()
    new RegisterApp()
})